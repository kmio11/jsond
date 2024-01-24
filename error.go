package jsond

import (
	"errors"
	"fmt"
)

var _ error = (*NodeError)(nil)

// NodeError represents an error that can occur during JSON processing in the jsond package.
type NodeError struct {
	code errCode
	path jsonpath
	err  error
}

type errCode int

const (
	codeUnknown = iota
	codeInternalError
	codeUnmarshalError
	codeMarshalError
	codeReadNullError
	codeReadUndefinedError
	codeSetNullError
	codeSetUndefinedError
)

func (e NodeError) Error() string {
	if len(e.path) > 0 {
		return fmt.Sprintf("%s at %s", e.err.Error(), e.path.String())
	}
	return e.err.Error()
}

func newInternalError(path jsonpath, err error) error {
	return &NodeError{
		code: codeInternalError,
		path: path,
		err:  fmt.Errorf("internal error. err=%w", err),
	}
}

// newReadNullError creates a new NodeError for attempting to read properties of null.
func newReadNullError(path jsonpath) error {
	prop := path[len(path)-1]

	return &NodeError{
		code: codeReadNullError,
		path: path,
		err:  fmt.Errorf("cannot read properties of null (reading '%v')", prop),
	}
}

// newReadUndefinedError creates a new NodeError for attempting to read properties of undefined.
func newReadUndefinedError(path jsonpath) error {
	prop := path[len(path)-1]

	return &NodeError{
		code: codeReadUndefinedError,
		path: path,
		err:  fmt.Errorf("cannot read properties of undefined (reading '%v')", prop),
	}
}

func newSetNullError(path jsonpath) error {
	prop := path[len(path)-1]

	return &NodeError{
		code: codeSetNullError,
		path: path,
		err:  fmt.Errorf("cannot set properties of null (setting '%v')", prop),
	}
}

func newSetUndefinedError(path jsonpath) error {
	prop := path[len(path)-1]

	return &NodeError{
		code: codeSetUndefinedError,
		path: path,
		err:  fmt.Errorf("cannot set properties of undefined (setting '%v')", prop),
	}
}

var _ error = (*Undefined)(nil)

// Undefined represents an undefined value.
type Undefined struct {
	path jsonpath
}

func (e Undefined) Error() string {
	return "undefined"
}

func newUndefined(path jsonpath) error {
	return &Undefined{
		path: path,
	}
}

// IsUndefined checks if the given error is Undefined.
func IsUndefined(err error) bool {
	target := &Undefined{}
	return errors.As(err, &target)
}
