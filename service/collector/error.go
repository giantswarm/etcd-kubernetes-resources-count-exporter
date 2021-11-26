package collector

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var emptyLineError = &microerror.Error{
	Kind: "emptyLineError",
}

// IsEmptyLine asserts emptyLineError.
func IsEmptyLine(err error) bool {
	return microerror.Cause(err) == emptyLineError
}
