package confetto

import "fmt"

// Param is the interface that all parameter types implement.
type Param interface {
	// key returns the configuration key for this parameter.
	key() string
	// setKey sets the configuration key.
	setKey(k string)
	// setFromString parses and sets the value from a string.
	setFromString(s string, listSeparator string) error
	// setFromAny sets the value from an arbitrary type (for YAML).
	setFromAny(v any, listSeparator string) error
	// validate runs all validators on the current value.
	validate() error
	// isRequired returns true if this parameter must be set.
	isRequired() bool
	// isSet returns true if the value has been explicitly set.
	IsSet() bool
	// hasDefault returns true if a default value was configured.
	hasDefault() bool
	// isSecret returns true if this parameter contains sensitive data.
	isSecret() bool
	// stringValue returns the current value formatted as a string.
	stringValue() string
}

// param is the internal generic parameter type that holds configuration for a single value.
type param[T any] struct {
	value      T
	defaultVal T
	hasDefVal  bool
	set        bool
	required   bool
	desc       string
	k          string
	secret     bool
	validators []func(T) error
}

func (p *param[T]) Get() T {
	return p.value
}

func (p *param[T]) IsSet() bool {
	return p.set
}

func (p *param[T]) key() string {
	return p.k
}

func (p *param[T]) setKey(k string) {
	p.k = k
}

func (p *param[T]) isRequired() bool {
	return p.required
}

func (p *param[T]) hasDefault() bool {
	return p.hasDefVal
}

func (p *param[T]) isSecret() bool {
	return p.secret
}

func (p *param[T]) stringValue() string {
	return fmt.Sprintf("%v", p.value)
}

func (p *param[T]) validate() error {
	for _, v := range p.validators {
		if err := v(p.value); err != nil {
			return &ValidationError{
				Key:     p.k,
				Value:   p.value,
				Message: err.Error(),
			}
		}
	}
	return nil
}
