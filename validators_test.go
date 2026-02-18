package confetto

import (
	"testing"
	"time"
)

func TestValidators(t *testing.T) {
	t.Run("Range", func(t *testing.T) {
		v := Range(1, 10)
		if err := v(5); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v(0); err == nil {
			t.Error("expected error for 0")
		}
		if err := v(11); err == nil {
			t.Error("expected error for 11")
		}
	})

	t.Run("OneOf", func(t *testing.T) {
		v := OneOf("a", "b", "c")
		if err := v("a"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v("d"); err == nil {
			t.Error("expected error for d")
		}
	})

	t.Run("MinLen", func(t *testing.T) {
		v := MinLen(3)
		if err := v("abc"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v("ab"); err == nil {
			t.Error("expected error for ab")
		}
	})

	t.Run("NotEmpty", func(t *testing.T) {
		v := NotEmpty()
		if err := v("a"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v(""); err == nil {
			t.Error("expected error for empty string")
		}
	})

	t.Run("Positive", func(t *testing.T) {
		v := Positive()
		if err := v(1); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v(0); err == nil {
			t.Error("expected error for 0")
		}
		if err := v(-1); err == nil {
			t.Error("expected error for -1")
		}
	})

	t.Run("NonNegative", func(t *testing.T) {
		v := NonNegative()
		if err := v(0); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v(-1); err == nil {
			t.Error("expected error for -1")
		}
	})

	t.Run("MinItems", func(t *testing.T) {
		v := MinItems[string](2)
		if err := v([]string{"a", "b"}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v([]string{"a"}); err == nil {
			t.Error("expected error for single item")
		}
	})

	t.Run("MaxItems", func(t *testing.T) {
		v := MaxItems[string](2)
		if err := v([]string{"a", "b"}); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if err := v([]string{"a", "b", "c"}); err == nil {
			t.Error("expected error for three items")
		}
	})
}

func TestValidators_MaxLen(t *testing.T) {
	v := MaxLen(5)

	t.Run("Valid", func(t *testing.T) {
		if err := v("abc"); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("Boundary", func(t *testing.T) {
		if err := v("abcde"); err != nil {
			t.Errorf("expected nil for exactly at limit, got %v", err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if err := v("abcdef"); err == nil {
			t.Error("expected error for string exceeding limit")
		}
	})
}

func TestValidators_RangeFloat(t *testing.T) {
	v := RangeFloat(1.0, 10.0)

	t.Run("Valid_InRange", func(t *testing.T) {
		if err := v(5.5); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("Valid_AtMin", func(t *testing.T) {
		if err := v(1.0); err != nil {
			t.Errorf("expected nil at min boundary, got %v", err)
		}
	})

	t.Run("Valid_AtMax", func(t *testing.T) {
		if err := v(10.0); err != nil {
			t.Errorf("expected nil at max boundary, got %v", err)
		}
	})

	t.Run("Invalid_BelowMin", func(t *testing.T) {
		if err := v(0.9); err == nil {
			t.Error("expected error for value below min")
		}
	})

	t.Run("Invalid_AboveMax", func(t *testing.T) {
		if err := v(10.1); err == nil {
			t.Error("expected error for value above max")
		}
	})
}

func TestValidators_RangeDuration(t *testing.T) {
	v := RangeDuration(time.Second, time.Minute)

	t.Run("Valid_InRange", func(t *testing.T) {
		if err := v(30 * time.Second); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("Valid_AtMin", func(t *testing.T) {
		if err := v(time.Second); err != nil {
			t.Errorf("expected nil at min boundary, got %v", err)
		}
	})

	t.Run("Valid_AtMax", func(t *testing.T) {
		if err := v(time.Minute); err != nil {
			t.Errorf("expected nil at max boundary, got %v", err)
		}
	})

	t.Run("Invalid_BelowMin", func(t *testing.T) {
		if err := v(500 * time.Millisecond); err == nil {
			t.Error("expected error for duration below min")
		}
	})

	t.Run("Invalid_AboveMax", func(t *testing.T) {
		if err := v(2 * time.Minute); err == nil {
			t.Error("expected error for duration above max")
		}
	})
}

func TestValidators_OneOfInt(t *testing.T) {
	v := OneOf(1, 2, 3)

	t.Run("Valid", func(t *testing.T) {
		if err := v(2); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if err := v(4); err == nil {
			t.Error("expected error for value not in set")
		}
	})
}
