package domain_test

import (
	"testing"
	"time"

	"venturex-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
)

var (
	now        = time.Now()
	pastTime   = now.Add(-1 * time.Hour)
	futureTime = now.Add(1 * time.Hour)
	zeroTime   = time.Time{}
)

func TestNewTimestamp(t *testing.T) {
	t.Run("NewTimestamp", func(t *testing.T) {
		got, err := domain.NewTimestamp(now)
		assert.NoError(t, err)
		assert.Equal(t, now.Unix(), got.Time().Unix())
		assert.Equal(t, now.Unix(), got.Value().Unix())
		assert.Equal(t, now.Unix(), got.Unix())
		assert.False(t, got.IsZero())

		_, err = domain.NewTimestamp(zeroTime)
		assert.ErrorIs(t, err, domain.ErrInvalidTime)
	})

	t.Run("NewTimestampNow", func(t *testing.T) {
		got := domain.NewTimestampNow()
		assert.False(t, got.IsZero())
		assert.WithinDuration(t, time.Now(), got.Time(), 1*time.Second)
	})

	t.Run("Timestamp methods", func(t *testing.T) {
		ts1, _ := domain.NewTimestamp(now)
		ts2, _ := domain.NewTimestamp(futureTime)
		assert.True(t, ts1.Before(ts2))
		assert.False(t, ts2.Before(ts1))
		assert.True(t, ts2.After(ts1))
		assert.False(t, ts1.After(ts2))
	})
}

func TestNewFutureTimestamp(t *testing.T) {
	_, err := domain.NewFutureTimestamp(zeroTime)
	assert.ErrorIs(t, err, domain.ErrInvalidTime)

	_, err = domain.NewFutureTimestamp(pastTime)
	assert.ErrorIs(t, err, domain.ErrPastTime)

	got, err := domain.NewFutureTimestamp(futureTime)
	assert.NoError(t, err)
	assert.Equal(t, futureTime.Unix(), got.Time().Unix())
	assert.Equal(t, futureTime.Unix(), got.Value().Unix())
	assert.Equal(t, futureTime.Unix(), got.Timestamp().Unix())
}

func TestNewPastTimestamp(t *testing.T) {
	_, err := domain.NewPastTimestamp(zeroTime)
	assert.ErrorIs(t, err, domain.ErrInvalidTime)

	_, err = domain.NewPastTimestamp(futureTime)
	assert.ErrorIs(t, err, domain.ErrFutureTime)

	got, err := domain.NewPastTimestamp(pastTime)
	assert.NoError(t, err)
	assert.Equal(t, pastTime.Unix(), got.Time().Unix())
	assert.Equal(t, pastTime.Unix(), got.Value().Unix())
	assert.Equal(t, pastTime.Unix(), got.Timestamp().Unix())
}

func TestNewCreatedAt(t *testing.T) {
	t.Run("NewCreatedAt", func(t *testing.T) {
		_, err := domain.NewCreatedAt(zeroTime)
		assert.ErrorIs(t, err, domain.ErrInvalidTime)

		_, err = domain.NewCreatedAt(now.Add(10 * time.Minute))
		assert.ErrorIs(t, err, domain.ErrFutureTime)

		got, err := domain.NewCreatedAt(now)
		assert.NoError(t, err)
		assert.Equal(t, now.Unix(), got.Time().Unix())
		assert.Equal(t, now.Unix(), got.Value().Unix())
		assert.Equal(t, now.Unix(), got.Timestamp().Unix())
	})

	t.Run("NewCreatedAtNow", func(t *testing.T) {
		got := domain.NewCreatedAtNow()
		assert.WithinDuration(t, time.Now(), got.Time(), 1*time.Second)
	})
}

func TestNewUpdatedAt(t *testing.T) {
	t.Run("NewUpdatedAt", func(t *testing.T) {
		_, err := domain.NewUpdatedAt(zeroTime)
		assert.ErrorIs(t, err, domain.ErrInvalidTime)

		_, err = domain.NewUpdatedAt(now.Add(10 * time.Minute))
		assert.ErrorIs(t, err, domain.ErrFutureTime)

		got, err := domain.NewUpdatedAt(now)
		assert.NoError(t, err)
		assert.Equal(t, now.Unix(), got.Time().Unix())
		assert.Equal(t, now.Unix(), got.Value().Unix())
		assert.Equal(t, now.Unix(), got.Timestamp().Unix())
	})

	t.Run("NewUpdatedAtNow", func(t *testing.T) {
		got := domain.NewUpdatedAtNow()
		assert.WithinDuration(t, time.Now(), got.Time(), 1*time.Second)
	})
}

func TestNewDuration(t *testing.T) {
	_, err := domain.NewDuration(-1 * time.Second)
	assert.ErrorIs(t, err, domain.ErrNegativeNumber)

	d := 10 * time.Second
	got, err := domain.NewDuration(d)
	assert.NoError(t, err)
	assert.Equal(t, d, got.Duration())
	assert.Equal(t, d, got.Value())
	assert.Equal(t, 10.0, got.Seconds())
	assert.InDelta(t, 10.0/60.0, got.Minutes(), 0.001)
	assert.InDelta(t, 10.0/3600.0, got.Hours(), 0.001)
}
