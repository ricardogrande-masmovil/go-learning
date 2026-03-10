package contextual_errors

import (
	"fmt"
	"testing"
)

func TestContextualError_InterfaceTrap(t *testing.T) {
	// A function that returns an error
	doFail := func(fail bool) error {
		var err error
		if fail {
			err = fmt.Errorf("underlying error")
		}
		// With(err) returns a *contextualError.
		// Err() ensures that if err was nil, untyped nil is returned natively.
		return With(err).Str("foo", "bar").Err()
	}

	err := doFail(false)
	if err != nil {
		t.Errorf("expected err to be nil, but got %v. Interface trap detected!", err)
	}

	err = doFail(true)
	if err == nil {
		t.Errorf("expected err to not be nil")
	}

	// Also check Wrap
	doFailWrap := func(fail bool) error {
		var err error
		if fail {
			err = fmt.Errorf("underlying error")
		}
		return Wrap(err, "wrapped").Str("foo", "bar").Err()
	}

	err = doFailWrap(false)
	if err != nil {
		t.Errorf("expected err to be nil, but got %v. Interface trap detected!", err)
	}
}
