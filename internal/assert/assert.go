package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("Got %v, but wanted %v", actual, expected)
	}
}

func Contains(t *testing.T, substring, container string) {
	t.Helper()
	if !strings.Contains(container, substring) {
		t.Errorf("Container %q doesnt contain substring %q", container, substring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
