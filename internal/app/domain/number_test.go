package domain_test

import (
	"testing"

	"beerdosan-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewPositiveInt(t *testing.T) {
	testCases := []struct {
		name      string
		value     int
		want      domain.PositiveInt
		expectErr error
	}{
		{"failure: zero", 0, 0, domain.ErrZeroValue},
		{"failure: negative", -10, 0, domain.ErrZeroValue},
		{"success", 100, 100, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewPositiveInt(tc.value)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, got.Value())
				assert.Equal(t, tc.value, got.Int())
			}
		})
	}
}

func TestNewNonNegativeInt(t *testing.T) {
	testCases := []struct {
		name      string
		value     int
		want      domain.NonNegativeInt
		expectErr error
	}{
		{"failure: negative", -1, 0, domain.ErrNegativeNumber},
		{"success: zero", 0, 0, nil},
		{"success: positive", 50, 50, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewNonNegativeInt(tc.value)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, got.Value())
				assert.Equal(t, tc.value == 0, got.IsZero())
			}
		})
	}
}

func TestNewAge(t *testing.T) {
	testCases := []struct {
		name      string
		value     int
		want      domain.Age
		expectErr error
	}{
		{"failure: negative", -1, 0, domain.ErrNegativeNumber},
		{"failure: too high", 151, 0, domain.ErrInvalidRange},
		{"success: lower bound", 0, 0, nil},
		{"success: upper bound", 150, 150, nil},
		{"success: valid age", 30, 30, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewAge(tc.value)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, got.Int())
			}
		})
	}
}

func TestNewPercentage(t *testing.T) {
	testCases := []struct {
		name      string
		value     int
		want      domain.Percentage
		expectErr error
	}{
		{"failure: negative", -1, 0, domain.ErrInvalidRange},
		{"failure: too high", 101, 0, domain.ErrInvalidRange},
		{"success: lower bound", 0, 0, nil},
		{"success: upper bound", 100, 100, nil},
		{"success: valid percentage", 50, 50, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewPercentage(tc.value)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, float64(tc.value)/100.0, got.Float64())
			}
		})
	}
}

func TestNewMoney(t *testing.T) {
	t.Run("from cents", func(t *testing.T) {
		testCases := []struct {
			name      string
			value     int64
			want      domain.Money
			expectErr error
		}{
			{"failure: negative", -1, 0, domain.ErrNegativeNumber},
			{"success: zero", 0, 0, nil},
			{"success: positive", 10050, 10050, nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := domain.NewMoney(tc.value)
				if tc.expectErr != nil {
					assert.ErrorIs(t, err, tc.expectErr)
					assert.Empty(t, got)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.want, got)
					assert.Equal(t, tc.value, got.Cents())
					assert.Equal(t, float64(tc.value)/100.0, got.Float64())
					assert.Equal(t, tc.value == 0, got.IsZero())
				}
			})
		}
	})

	t.Run("from float", func(t *testing.T) {
		testCases := []struct {
			name      string
			value     float64
			want      domain.Money
			expectErr error
		}{
			{"failure: negative", -0.01, 0, domain.ErrNegativeNumber},
			{"success: zero", 0.0, 0, nil},
			{"success: positive", 123.45, 12345, nil},
			{"success: rounding", 123.456, 12345, nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := domain.NewMoneyFromFloat(tc.value)
				if tc.expectErr != nil {
					assert.ErrorIs(t, err, tc.expectErr)
					assert.Empty(t, got)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.want, got)
				}
			})
		}
	})
}

func TestNewID(t *testing.T) {
	testCases := []struct {
		name      string
		value     int64
		want      domain.ID
		expectErr error
	}{
		{"failure: zero", 0, 0, domain.ErrZeroValue},
		{"failure: negative", -1, 0, domain.ErrZeroValue},
		{"success", 12345, 12345, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewID(tc.value)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, got.Value())
				assert.Equal(t, tc.value, got.Int64())
				assert.Equal(t, int(tc.value), got.Int())
			}
		})
	}
}
