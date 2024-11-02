package flags

import "github.com/pasataleo/go-errors/errors"

const (
	ErrorCodeMissingFlag   errors.ErrorCode = "flags.ErrorCodeMissingArg"
	ErrorCodeUnknownFlag   errors.ErrorCode = "flags.ErrorCodeUnknownFlag"
	ErrorCodeDuplicateFlag errors.ErrorCode = "flags.ErrorCodeDuplicateFlag"
	ErrorCodeInvalidValue  errors.ErrorCode = "flags.ErrorCodeInvalidValue"
)
