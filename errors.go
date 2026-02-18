package confetto

import (
	"fmt"
	"strings"
)

// LoadError contains all errors that occurred during configuration loading.
type LoadError struct {
	Errors []error
}

func (e *LoadError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d configuration errors:\n", len(e.Errors)))
	for _, err := range e.Errors {
		sb.WriteString("  - ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e *LoadError) Add(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *LoadError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ParseError indicates that a value could not be parsed to the expected type.
type ParseError struct {
	Key      string
	Value    string
	Expected string
	Err      error
}

func (e *ParseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(
			"failed to parse %q as %s for key %q: %v", e.Value, e.Expected, e.Key, e.Err,
		)
	}
	return fmt.Sprintf("failed to parse %q as %s for key %q", e.Value, e.Expected, e.Key)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// ValidationError indicates that a value failed validation.
type ValidationError struct {
	Key     string
	Value   any
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %q (value: %v): %s", e.Key, e.Value, e.Message)
}

// RequiredError indicates that a required parameter was not set.
type RequiredError struct {
	Key string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("required parameter %q is not set", e.Key)
}
