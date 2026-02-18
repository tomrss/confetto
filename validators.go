package confetto

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// ErrValidation is the sentinel error for validation failures.
var ErrValidation = errors.New("validation error")

// Range returns a validator that checks if an int is within [min, max].
func Range(lo, hi int) func(int) error {
	return func(v int) error {
		if v < lo || v > hi {
			return fmt.Errorf("%w: value %d is not in range [%d, %d]", ErrValidation, v, lo, hi)
		}
		return nil
	}
}

// RangeFloat returns a validator that checks if a float64 is within [min, max].
func RangeFloat(lo, hi float64) func(float64) error {
	return func(v float64) error {
		if v < lo || v > hi {
			return fmt.Errorf("%w: value %f is not in range [%f, %f]", ErrValidation, v, lo, hi)
		}
		return nil
	}
}

// RangeDuration returns a validator that checks if a duration is within [min, max].
func RangeDuration(lo, hi time.Duration) func(time.Duration) error {
	return func(v time.Duration) error {
		if v < lo || v > hi {
			return fmt.Errorf("%w: value %v is not in range [%v, %v]", ErrValidation, v, lo, hi)
		}
		return nil
	}
}

// OneOf returns a validator that checks if a value is one of the allowed values.
func OneOf[T comparable](allowed ...T) func(T) error {
	return func(v T) error {
		if slices.Contains(allowed, v) {
			return nil
		}
		return fmt.Errorf("%w: value %v is not one of %v", ErrValidation, v, allowed)
	}
}

// MinLen returns a validator that checks if a string has at least n characters.
func MinLen(n int) func(string) error {
	return func(v string) error {
		if len(v) < n {
			return fmt.Errorf(
				"%w: string length %d is less than minimum %d", ErrValidation, len(v), n,
			)
		}
		return nil
	}
}

// MaxLen returns a validator that checks if a string has at most n characters.
func MaxLen(n int) func(string) error {
	return func(v string) error {
		if len(v) > n {
			return fmt.Errorf(
				"%w: string length %d is greater than maximum %d", ErrValidation, len(v), n,
			)
		}
		return nil
	}
}

// NotEmpty returns a validator that checks if a string is not empty.
func NotEmpty() func(string) error {
	return func(v string) error {
		if v == "" {
			return fmt.Errorf("%w: string must not be empty", ErrValidation)
		}
		return nil
	}
}

// MinItems returns a validator that checks if a slice has at least n items.
func MinItems[T any](n int) func([]T) error {
	return func(v []T) error {
		if len(v) < n {
			return fmt.Errorf("%w: list has %d items, minimum is %d", ErrValidation, len(v), n)
		}
		return nil
	}
}

// MaxItems returns a validator that checks if a slice has at most n items.
func MaxItems[T any](n int) func([]T) error {
	return func(v []T) error {
		if len(v) > n {
			return fmt.Errorf("%w: list has %d items, maximum is %d", ErrValidation, len(v), n)
		}
		return nil
	}
}

// Positive returns a validator that checks if an int is positive (> 0).
func Positive() func(int) error {
	return func(v int) error {
		if v <= 0 {
			return fmt.Errorf("%w: value %d must be positive", ErrValidation, v)
		}
		return nil
	}
}

// NonNegative returns a validator that checks if an int is non-negative (>= 0).
func NonNegative() func(int) error {
	return func(v int) error {
		if v < 0 {
			return fmt.Errorf("%w: value %d must be non-negative", ErrValidation, v)
		}
		return nil
	}
}
