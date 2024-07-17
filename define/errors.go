package define

import "fmt"

type UndefinedTypeError string

func (e UndefinedTypeError) Error() string {
	return fmt.Sprintf("undefined type: %s", string(e))
}

type TypeAssertionError struct {
	source, target string
}

func (e TypeAssertionError) Error() string {
	return fmt.Sprintf("type assertion failed: expected [%s], given [%s]", e.target, e.source)
}

func (e TypeAssertionError) SourceType() string {
	return e.source
}

func (e TypeAssertionError) TargetType() string {
	return e.target
}
