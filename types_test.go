package confetto

import "testing"

func TestLoad_ParamGetBeforeSet(t *testing.T) {
	t.Run("StringParam", func(t *testing.T) {
		p := String().Build()
		if p.Get() != "" {
			t.Errorf("expected empty string, got %q", p.Get())
		}
	})

	t.Run("IntParam", func(t *testing.T) {
		p := Int().Build()
		if p.Get() != 0 {
			t.Errorf("expected 0, got %d", p.Get())
		}
	})

	t.Run("BoolParam", func(t *testing.T) {
		p := Bool().Build()
		if p.Get() != false {
			t.Errorf("expected false, got %v", p.Get())
		}
	})

	t.Run("FloatParam", func(t *testing.T) {
		p := Float().Build()
		if p.Get() != 0.0 {
			t.Errorf("expected 0.0, got %f", p.Get())
		}
	})

	t.Run("DurationParam", func(t *testing.T) {
		p := Duration().Build()
		if p.Get() != 0 {
			t.Errorf("expected 0, got %v", p.Get())
		}
	})

	t.Run("StringListParam", func(t *testing.T) {
		p := StringList().Build()
		if p.Get() != nil {
			t.Errorf("expected nil, got %v", p.Get())
		}
	})
}

func TestLoad_IsSetBehavior(t *testing.T) {
	type isSetConfig struct {
		Name StringParam `cfg:"name"`
	}

	cfg := isSetConfig{
		Name: String().Default("default").Build(),
	}

	// before Load, IsSet should be false
	if cfg.Name.IsSet() {
		t.Error("expected IsSet() == false before Load")
	}

	err := Load(&cfg, Options{Args: []string{"--name=hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// after Load with explicit value, IsSet should be true
	if !cfg.Name.IsSet() {
		t.Error("expected IsSet() == true after Load with explicit value")
	}
}
