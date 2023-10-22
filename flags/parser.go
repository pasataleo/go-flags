package flags

import "github.com/pasataleo/go-errors/errors"

type Parser[T any] interface {
	Parse(name string, args map[string]string) (T, error)
}

type StringParser struct{}

func (p *StringParser) Parse(name string, args map[string]string) (string, error) {
	if len(args) == 0 {
		return "", errors.Newf(nil, ErrorCodeMissingFlag, "Missing flag %s", name)
	}

	if len(args) > 1 {
		return "", errors.Newf(nil, ErrorCodeDuplicateFlag, "Duplicate flag %s", name)
	}

	return args[name], nil
}
