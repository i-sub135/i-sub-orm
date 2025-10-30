package constant

import "errors"

var (
	ErrDestination     = errors.New("destination must be pointer")
	ErrDestinationType = errors.New("destination must be slice or struct")
)
