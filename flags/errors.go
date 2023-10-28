package flags

import "github.com/pasataleo/go-errors/errors"

const (
	ErrorCodeMissingFlag   errors.ErrorCode = "Flags.ErrorCodeMissingArg"
	ErrorCodeUnknownFlag   errors.ErrorCode = "Flags.ErrorCodeUnknownFlag"
	ErrorCodeDuplicateFlag errors.ErrorCode = "Flags.ErrorCodeDuplicateFlag"
	ErrorCodeInvalidValue  errors.ErrorCode = "Flags.ErrorCodeInvalidValue"
)
