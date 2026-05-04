package enums

import "fmt"

// validationError is returned by every enum's Validate() method when the
// receiver carries a value not in the canonical set. It implements the
// error interface and is comparable with errors.Is against
// [ErrInvalidEnum].
type validationError struct {
	enum  string
	value string
}

func (e *validationError) Error() string {
	return fmt.Sprintf("enums: invalid %s value %q", e.enum, e.value)
}

// Is satisfies errors.Is so callers can match without unwrapping. Every
// enum-validation failure unwraps to [ErrInvalidEnum].
func (e *validationError) Is(target error) bool { return target == ErrInvalidEnum }

// ErrInvalidEnum is the sentinel returned from every enum's Validate
// method when the receiver isn't one of the defined wire values. Use
// errors.Is to detect it.
var ErrInvalidEnum = &validationError{enum: "<unknown>", value: ""}

// invalid is the package-internal helper each Validate method calls to
// build a concrete error if Valid() returned false.
func invalid(enum, value string) error {
	return &validationError{enum: enum, value: value}
}
