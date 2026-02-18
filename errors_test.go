package confetto

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestLoad_ErrorMessages(t *testing.T) {
	t.Run("LoadError_SingleError", func(t *testing.T) {
		le := &LoadError{}
		le.Add(fmt.Errorf("single error"))
		msg := le.Error()
		if msg != "single error" {
			t.Errorf("expected 'single error', got %q", msg)
		}
	})

	t.Run("LoadError_MultipleErrors", func(t *testing.T) {
		le := &LoadError{}
		le.Add(fmt.Errorf("error one"))
		le.Add(fmt.Errorf("error two"))
		msg := le.Error()
		if !strings.Contains(msg, "2 configuration errors") {
			t.Errorf("expected '2 configuration errors' header, got %q", msg)
		}
		if !strings.Contains(msg, "error one") || !strings.Contains(msg, "error two") {
			t.Errorf("expected both errors in message, got %q", msg)
		}
	})

	t.Run("ParseError_NilErr", func(t *testing.T) {
		pe := &ParseError{Key: "mykey", Value: "myval", Expected: "int"}
		msg := pe.Error()
		if !strings.Contains(msg, "mykey") || !strings.Contains(msg, "myval") {
			t.Errorf("expected key and value in message, got %q", msg)
		}
		// should not contain the wrapping error text
		if strings.Contains(msg, ": <nil>") {
			t.Errorf("should not contain nil error text, got %q", msg)
		}
	})

	t.Run("ParseError_Unwrap", func(t *testing.T) {
		inner := fmt.Errorf("inner error")
		pe := &ParseError{Key: "k", Value: "v", Expected: "int", Err: inner}
		if !errors.Is(pe, inner) {
			t.Error("expected Unwrap to return inner error")
		}

		peNil := &ParseError{Key: "k", Value: "v", Expected: "int"}
		if peNil.Unwrap() != nil {
			t.Error("expected Unwrap to return nil when Err is nil")
		}
	})

	t.Run("RequiredError", func(t *testing.T) {
		re := &RequiredError{Key: "db.host"}
		msg := re.Error()
		if !strings.Contains(msg, "db.host") {
			t.Errorf("expected key in message, got %q", msg)
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		ve := &ValidationError{Key: "port", Value: 99999, Message: "out of range"}
		msg := ve.Error()
		if !strings.Contains(msg, "port") || !strings.Contains(msg, "99999") || !strings.Contains(msg, "out of range") {
			t.Errorf("expected key, value and message in error, got %q", msg)
		}
	})
}
