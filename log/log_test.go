// Licensed under Elastic License 2.0
// See LICENSE.txt for details

package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

func TestLoggerShouldntPrint(t *testing.T) {
	var b bytes.Buffer
	e := errors.New("test error happened")
	l := New(&b)
	l.WithTrace("1234").Debugf("error: %s", e)
	if len(b.String()) > 0 {
		t.Fatalf("expected empty string, got %s", b.String())
	}
}

func TestLoggerTrace(t *testing.T) {
	var b bytes.Buffer
	e := errors.New("test error happened")
	l := New(&b)
	l.WithTrace("1234").Errorf("error: %s", e)
	rec := &Record{}
	err := json.Unmarshal(b.Bytes(), rec)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Message != "error: test error happened" {
		t.Fatalf("expected error message, got %s", rec.Message)
	}

	if rec.TraceID != "1234" {
		t.Fatalf("expected trace '1234', got %s", rec.TraceID)
	}
}

func TestLoggerNoTrace(t *testing.T) {
	var b bytes.Buffer
	e := errors.New("test error happened")
	l := New(&b)
	l.Errorf("error: %s", e)

	rec := &Record{}
	err := json.Unmarshal(b.Bytes(), rec)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Message != "error: test error happened" {
		t.Fatalf("expected error message, got %s", rec.Message)
	}

	if rec.TraceID != "" {
		t.Fatalf("expected no traceID, got %s", rec.TraceID)
	}
}

func TestDefaultLogger(t *testing.T) {
	e := errors.New("test error happened")
	WithTrace("1234").Errorf("error: %s", e)
}

func TestDefaultWith(t *testing.T) {
	e := errors.New("test error happened")
	x := WithTrace("1234").With("newfield", "my new field")
	x.Errorf("oh no %s", e)
}

func TestPrintf(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	l.Printf("hello %s", "user")
	if b.String() != "hello user\n" {
		t.Fatalf("expected plain string, got %s", b.String())
	}
}

func TestPrintToScreen(t *testing.T) {
	Printf("hello %s", "user")
}
