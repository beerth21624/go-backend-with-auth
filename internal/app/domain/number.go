package domain

import "errors"

var (
	ErrNegativeNumber = errors.New("number cannot be negative")
	ErrZeroValue      = errors.New("value cannot be zero")
	ErrInvalidRange   = errors.New("value out of valid range")
)

type PositiveInt int

func NewPositiveInt(n int) (PositiveInt, error) {
	if n <= 0 {
		return 0, ErrZeroValue
	}
	return PositiveInt(n), nil
}

func (p PositiveInt) Value() int {
	return int(p)
}

func (p PositiveInt) Int() int {
	return int(p)
}

type NonNegativeInt int

func NewNonNegativeInt(n int) (NonNegativeInt, error) {
	if n < 0 {
		return 0, ErrNegativeNumber
	}
	return NonNegativeInt(n), nil
}

func (n NonNegativeInt) Value() int {
	return int(n)
}

func (n NonNegativeInt) Int() int {
	return int(n)
}

func (n NonNegativeInt) IsZero() bool {
	return n == 0
}

type Age int

func NewAge(n int) (Age, error) {
	if n < 0 {
		return 0, ErrNegativeNumber
	}
	if n > 150 {
		return 0, ErrInvalidRange
	}
	return Age(n), nil
}

func (a Age) Value() int {
	return int(a)
}

func (a Age) Int() int {
	return int(a)
}

type Percentage int

func NewPercentage(n int) (Percentage, error) {
	if n < 0 || n > 100 {
		return 0, ErrInvalidRange
	}
	return Percentage(n), nil
}

func (p Percentage) Value() int {
	return int(p)
}

func (p Percentage) Int() int {
	return int(p)
}

func (p Percentage) Float64() float64 {
	return float64(p) / 100.0
}

type Money int64

func NewMoney(cents int64) (Money, error) {
	if cents < 0 {
		return 0, ErrNegativeNumber
	}
	return Money(cents), nil
}

func NewMoneyFromFloat(amount float64) (Money, error) {
	cents := int64(amount * 100)
	return NewMoney(cents)
}

func (m Money) Cents() int64 {
	return int64(m)
}

func (m Money) Float64() float64 {
	return float64(m) / 100.0
}

func (m Money) IsZero() bool {
	return m == 0
}

type ID int64

func NewID(n int64) (ID, error) {
	if n <= 0 {
		return 0, ErrZeroValue
	}
	return ID(n), nil
}

func (id ID) Value() int64 {
	return int64(id)
}

func (id ID) Int64() int64 {
	return int64(id)
}

func (id ID) Int() int {
	return int(id)
}
