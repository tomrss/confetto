package confetto

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// StringParam holds a string configuration value.
type StringParam struct {
	param[string]
}

//nolint:unparam // error is always nil but signature must match other param types
func (p *StringParam) setFromString(s string, _ string) error {
	p.value = s
	p.set = true
	return nil
}

//nolint:unparam // error is always nil but signature must match other param types
func (p *StringParam) setFromAny(v any, _ string) error {
	switch val := v.(type) {
	case string:
		p.value = val
	default:
		p.value = fmt.Sprintf("%v", val)
	}
	p.set = true
	return nil
}

// IntParam holds an int configuration value.
type IntParam struct {
	param[int]
}

func (p *IntParam) setFromString(s string, _ string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return &ParseError{Key: p.k, Value: s, Expected: "int", Err: err}
	}
	p.value = v
	p.set = true
	return nil
}

func (p *IntParam) setFromAny(v any, _ string) error {
	switch val := v.(type) {
	case int:
		p.value = val
	case int64:
		p.value = int(val)
	case float64:
		p.value = int(val)
	case string:
		return p.setFromString(val, "")
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "int"}
	}
	p.set = true
	return nil
}

// BoolParam holds a bool configuration value.
type BoolParam struct {
	param[bool]
}

func (p *BoolParam) setFromString(s string, _ string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return &ParseError{Key: p.k, Value: s, Expected: "bool", Err: err}
	}
	p.value = v
	p.set = true
	return nil
}

func (p *BoolParam) setFromAny(v any, _ string) error {
	switch val := v.(type) {
	case bool:
		p.value = val
	case string:
		return p.setFromString(val, "")
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "bool"}
	}
	p.set = true
	return nil
}

// FloatParam holds a float64 configuration value.
type FloatParam struct {
	param[float64]
}

func (p *FloatParam) setFromString(s string, _ string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return &ParseError{Key: p.k, Value: s, Expected: "float64", Err: err}
	}
	p.value = v
	p.set = true
	return nil
}

func (p *FloatParam) setFromAny(v any, _ string) error {
	switch val := v.(type) {
	case float64:
		p.value = val
	case float32:
		p.value = float64(val)
	case int:
		p.value = float64(val)
	case int64:
		p.value = float64(val)
	case string:
		return p.setFromString(val, "")
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "float64"}
	}
	p.set = true
	return nil
}

// DurationParam holds a time.Duration configuration value.
type DurationParam struct {
	param[time.Duration]
}

func (p *DurationParam) setFromString(s string, _ string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return &ParseError{Key: p.k, Value: s, Expected: "duration", Err: err}
	}
	p.value = v
	p.set = true
	return nil
}

func (p *DurationParam) setFromAny(v any, _ string) error {
	switch val := v.(type) {
	case string:
		return p.setFromString(val, "")
	case int:
		p.value = time.Duration(val)
	case int64:
		p.value = time.Duration(val)
	case float64:
		p.value = time.Duration(val)
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "duration"}
	}
	p.set = true
	return nil
}

// StringListParam holds a []string configuration value.
type StringListParam struct {
	param[[]string]
}

func (p *StringListParam) setFromString(s string, sep string) error {
	if s == "" {
		p.value = []string{}
	} else {
		p.value = strings.Split(s, sep)
	}
	p.set = true
	return nil
}

func (p *StringListParam) setFromAny(v any, sep string) error {
	switch val := v.(type) {
	case []any:
		p.value = make([]string, len(val))
		for i, item := range val {
			p.value[i] = fmt.Sprintf("%v", item)
		}
	case []string:
		p.value = val
	case string:
		return p.setFromString(val, sep)
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "[]string"}
	}
	p.set = true
	return nil
}

// IntListParam holds a []int configuration value.
type IntListParam struct {
	param[[]int]
}

func (p *IntListParam) setFromString(s string, sep string) error {
	if s == "" {
		p.value = []int{}
		p.set = true
		return nil
	}
	parts := strings.Split(s, sep)
	p.value = make([]int, len(parts))
	for i, part := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return &ParseError{Key: p.k, Value: part, Expected: "int", Err: err}
		}
		p.value[i] = v
	}
	p.set = true
	return nil
}

func (p *IntListParam) setFromAny(v any, sep string) error {
	switch val := v.(type) {
	case []any:
		p.value = make([]int, len(val))
		for i, item := range val {
			switch n := item.(type) {
			case int:
				p.value[i] = n
			case int64:
				p.value[i] = int(n)
			case float64:
				p.value[i] = int(n)
			default:
				return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", item), Expected: "int"}
			}
		}
	case []int:
		p.value = val
	case string:
		return p.setFromString(val, sep)
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "[]int"}
	}
	p.set = true
	return nil
}

// BoolListParam holds a []bool configuration value.
type BoolListParam struct {
	param[[]bool]
}

func (p *BoolListParam) setFromString(s string, sep string) error {
	if s == "" {
		p.value = []bool{}
		p.set = true
		return nil
	}
	parts := strings.Split(s, sep)
	p.value = make([]bool, len(parts))
	for i, part := range parts {
		v, err := strconv.ParseBool(strings.TrimSpace(part))
		if err != nil {
			return &ParseError{Key: p.k, Value: part, Expected: "bool", Err: err}
		}
		p.value[i] = v
	}
	p.set = true
	return nil
}

func (p *BoolListParam) setFromAny(v any, sep string) error {
	switch val := v.(type) {
	case []any:
		p.value = make([]bool, len(val))
		for i, item := range val {
			switch b := item.(type) {
			case bool:
				p.value[i] = b
			default:
				return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", item), Expected: "bool"}
			}
		}
	case []bool:
		p.value = val
	case string:
		return p.setFromString(val, sep)
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "[]bool"}
	}
	p.set = true
	return nil
}

// FloatListParam holds a []float64 configuration value.
type FloatListParam struct {
	param[[]float64]
}

func (p *FloatListParam) setFromString(s string, sep string) error {
	if s == "" {
		p.value = []float64{}
		p.set = true
		return nil
	}
	parts := strings.Split(s, sep)
	p.value = make([]float64, len(parts))
	for i, part := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return &ParseError{Key: p.k, Value: part, Expected: "float64", Err: err}
		}
		p.value[i] = v
	}
	p.set = true
	return nil
}

func (p *FloatListParam) setFromAny(v any, sep string) error {
	switch val := v.(type) {
	case []any:
		p.value = make([]float64, len(val))
		for i, item := range val {
			switch n := item.(type) {
			case float64:
				p.value[i] = n
			case float32:
				p.value[i] = float64(n)
			case int:
				p.value[i] = float64(n)
			case int64:
				p.value[i] = float64(n)
			default:
				return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", item), Expected: "float64"}
			}
		}
	case []float64:
		p.value = val
	case string:
		return p.setFromString(val, sep)
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "[]float64"}
	}
	p.set = true
	return nil
}

// DurationListParam holds a []time.Duration configuration value.
type DurationListParam struct {
	param[[]time.Duration]
}

func (p *DurationListParam) setFromString(s string, sep string) error {
	if s == "" {
		p.value = []time.Duration{}
		p.set = true
		return nil
	}
	parts := strings.Split(s, sep)
	p.value = make([]time.Duration, len(parts))
	for i, part := range parts {
		v, err := time.ParseDuration(strings.TrimSpace(part))
		if err != nil {
			return &ParseError{Key: p.k, Value: part, Expected: "duration", Err: err}
		}
		p.value[i] = v
	}
	p.set = true
	return nil
}

func (p *DurationListParam) setFromAny(v any, sep string) error {
	switch val := v.(type) {
	case []any:
		p.value = make([]time.Duration, len(val))
		for i, item := range val {
			switch d := item.(type) {
			case string:
				parsed, err := time.ParseDuration(d)
				if err != nil {
					return &ParseError{Key: p.k, Value: d, Expected: "duration", Err: err}
				}
				p.value[i] = parsed
			case int:
				p.value[i] = time.Duration(d)
			case int64:
				p.value[i] = time.Duration(d)
			case float64:
				p.value[i] = time.Duration(d)
			default:
				return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", item), Expected: "duration"}
			}
		}
	case []time.Duration:
		p.value = val
	case string:
		return p.setFromString(val, sep)
	default:
		return &ParseError{Key: p.k, Value: fmt.Sprintf("%v", v), Expected: "[]duration"}
	}
	p.set = true
	return nil
}
