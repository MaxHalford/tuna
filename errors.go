package tuna

import "fmt"

// An ErrUnknownField occurs trying to access an unexisting Row field.
type ErrUnknownField struct {
	field string
}

// Error implements the Error interface.
func (e ErrUnknownField) Error() string {
	return fmt.Sprintf("no field named %s", e.field)
}
