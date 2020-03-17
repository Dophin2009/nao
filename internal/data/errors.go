package data

import "errors"

var (
	// errNil is an error returned when some pointer is nil.
	errNil = errors.New("is nil")
	// errInvalid is an error returned when some value is invalid.
	errInvalid = errors.New("invalid")
	// errAlreadyExists is an error returned when a unique value already exists.
	errAlreadyExists = errors.New("already exists")
)

const (
	errmsgModelAssertType = "failed to assert type of model"

	errmsgJSONMarshal   = "failed to marshal to JSON"
	errmsgJSONUnmarshal = "failed to unmarshal from JSON"
)
