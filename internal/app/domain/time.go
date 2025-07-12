package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidTime = errors.New("invalid time")
	ErrPastTime    = errors.New("time cannot be in the past")
	ErrFutureTime  = errors.New("time cannot be in the future")
)

type Timestamp time.Time

func NewTimestamp(t time.Time) (Timestamp, error) {
	if t.IsZero() {
		return Timestamp{}, ErrInvalidTime
	}
	return Timestamp(t), nil
}

func NewTimestampNow() Timestamp {
	return Timestamp(time.Now())
}

func (ts Timestamp) Time() time.Time {
	return time.Time(ts)
}

func (ts Timestamp) Value() time.Time {
	return time.Time(ts)
}

func (ts Timestamp) Unix() int64 {
	return time.Time(ts).Unix()
}

func (ts Timestamp) IsZero() bool {
	return time.Time(ts).IsZero()
}

func (ts Timestamp) Before(other Timestamp) bool {
	return time.Time(ts).Before(time.Time(other))
}

func (ts Timestamp) After(other Timestamp) bool {
	return time.Time(ts).After(time.Time(other))
}

type FutureTimestamp time.Time

func NewFutureTimestamp(t time.Time) (FutureTimestamp, error) {
	if t.IsZero() {
		return FutureTimestamp{}, ErrInvalidTime
	}
	if t.Before(time.Now()) {
		return FutureTimestamp{}, ErrPastTime
	}
	return FutureTimestamp(t), nil
}

func (ft FutureTimestamp) Time() time.Time {
	return time.Time(ft)
}

func (ft FutureTimestamp) Value() time.Time {
	return time.Time(ft)
}

func (ft FutureTimestamp) Timestamp() Timestamp {
	return Timestamp(ft)
}

type PastTimestamp time.Time

func NewPastTimestamp(t time.Time) (PastTimestamp, error) {
	if t.IsZero() {
		return PastTimestamp{}, ErrInvalidTime
	}
	if t.After(time.Now()) {
		return PastTimestamp{}, ErrFutureTime
	}
	return PastTimestamp(t), nil
}

func (pt PastTimestamp) Time() time.Time {
	return time.Time(pt)
}

func (pt PastTimestamp) Value() time.Time {
	return time.Time(pt)
}

func (pt PastTimestamp) Timestamp() Timestamp {
	return Timestamp(pt)
}

type CreatedAt time.Time

func NewCreatedAt(t time.Time) (CreatedAt, error) {
	if t.IsZero() {
		return CreatedAt{}, ErrInvalidTime
	}
	if t.After(time.Now().Add(5 * time.Minute)) {
		return CreatedAt{}, ErrFutureTime
	}
	return CreatedAt(t), nil
}

func NewCreatedAtNow() CreatedAt {
	return CreatedAt(time.Now())
}

func (ca CreatedAt) Time() time.Time {
	return time.Time(ca)
}

func (ca CreatedAt) Value() time.Time {
	return time.Time(ca)
}

func (ca CreatedAt) Timestamp() Timestamp {
	return Timestamp(ca)
}

type UpdatedAt time.Time

func NewUpdatedAt(t time.Time) (UpdatedAt, error) {
	if t.IsZero() {
		return UpdatedAt{}, ErrInvalidTime
	}
	if t.After(time.Now().Add(5 * time.Minute)) {
		return UpdatedAt{}, ErrFutureTime
	}
	return UpdatedAt(t), nil
}

func NewUpdatedAtNow() UpdatedAt {
	return UpdatedAt(time.Now())
}

func (ua UpdatedAt) Time() time.Time {
	return time.Time(ua)
}

func (ua UpdatedAt) Value() time.Time {
	return time.Time(ua)
}

func (ua UpdatedAt) Timestamp() Timestamp {
	return Timestamp(ua)
}

type Duration time.Duration

func NewDuration(d time.Duration) (Duration, error) {
	if d < 0 {
		return 0, ErrNegativeNumber
	}
	return Duration(d), nil
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) Value() time.Duration {
	return time.Duration(d)
}

func (d Duration) Seconds() float64 {
	return time.Duration(d).Seconds()
}

func (d Duration) Minutes() float64 {
	return time.Duration(d).Minutes()
}

func (d Duration) Hours() float64 {
	return time.Duration(d).Hours()
}
