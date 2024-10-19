package flags

import (
	"strconv"
	"strings"

	"github.com/pasataleo/go-errors/errors"
)

type Parser[T any] interface {
	Parse(name string, args []string) (T, error)
}

var _ Parser[any] = ParserFn[any](nil)

type ParserFn[T any] func(name string, args []string) (T, error)

func (fn ParserFn[T]) Parse(name string, args []string) (T, error) {
	return fn(name, args)
}

type aliasParser[T any] interface {
	Parse(name string, args map[string][]string) (T, error)
}

type parserWrapper[T any] struct {
	parser Parser[T]
}

func (p *parserWrapper[T]) Parse(name string, args []string) (interface{}, error) {
	return p.parser.Parse(name, args)
}

type aliasWrapper[T any] struct {
	parser aliasParser[T]
}

func (p *aliasWrapper[T]) Parse(name string, args map[string][]string) (interface{}, error) {
	return p.parser.Parse(name, args)
}

type singleArgParser[T any] struct {
	parser func(arg string) (T, error)
}

func (p *singleArgParser[T]) Parse(name string, args []string) (T, error) {
	var errorResult T

	if len(args) == 0 {
		return errorResult, errors.Newf(nil, ErrorCodeMissingFlag, "missing flag %q", name)
	}

	if len(args) > 1 {
		return errorResult, errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %q", name)
	}

	value, err := p.parser(args[0])
	if err != nil {
		return errorResult, errors.Newf(err, ErrorCodeInvalidValue, "invalid value for flag %q", name)
	}
	return value, nil
}

type sliceArgParser[T any] struct {
	parser func(arg string) (T, error)
}

func (p *sliceArgParser[T]) Parse(name string, args []string) ([]T, error) {
	if len(args) == 0 {
		return nil, errors.Newf(nil, ErrorCodeMissingFlag, "missing flag %q", name)
	}

	var result []T
	for _, arg := range args {
		value, err := p.parser(arg)
		if err != nil {
			return nil, errors.Newf(err, ErrorCodeInvalidValue, "invalid value for flag %q", name)
		}
		result = append(result, value)
	}
	return result, nil
}

func intParser() Parser[int] {
	return &singleArgParser[int]{
		parser: strconv.Atoi,
	}
}

func intSliceParser() Parser[[]int] {
	return &sliceArgParser[int]{
		parser: strconv.Atoi,
	}
}

func int8Parser() Parser[int8] {
	return &singleArgParser[int8]{
		parser: func(arg string) (int8, error) {
			value, err := strconv.ParseInt(arg, 10, 8)
			return int8(value), err
		},
	}
}

func int8SliceParser() Parser[[]int8] {
	return &sliceArgParser[int8]{
		parser: func(arg string) (int8, error) {
			value, err := strconv.ParseInt(arg, 10, 8)
			return int8(value), err
		},
	}
}

func int16Parser() Parser[int16] {
	return &singleArgParser[int16]{
		parser: func(arg string) (int16, error) {
			value, err := strconv.ParseInt(arg, 10, 16)
			return int16(value), err
		},
	}
}

func int16SliceParser() Parser[[]int16] {
	return &sliceArgParser[int16]{
		parser: func(arg string) (int16, error) {
			value, err := strconv.ParseInt(arg, 10, 16)
			return int16(value), err
		},
	}
}

func int32Parser() Parser[int32] {
	return &singleArgParser[int32]{
		parser: func(arg string) (int32, error) {
			value, err := strconv.ParseInt(arg, 10, 32)
			return int32(value), err
		},
	}
}

func int32SliceParser() Parser[[]int32] {
	return &sliceArgParser[int32]{
		parser: func(arg string) (int32, error) {
			value, err := strconv.ParseInt(arg, 10, 32)
			return int32(value), err
		},
	}
}

func int64Parser() Parser[int64] {
	return &singleArgParser[int64]{
		parser: func(arg string) (int64, error) {
			value, err := strconv.ParseInt(arg, 10, 64)
			return int64(value), err
		},
	}
}

func int64SliceParser() Parser[[]int64] {
	return &sliceArgParser[int64]{
		parser: func(arg string) (int64, error) {
			value, err := strconv.ParseInt(arg, 10, 64)
			return int64(value), err
		},
	}
}

func uintParser() Parser[uint] {
	return &singleArgParser[uint]{
		parser: func(arg string) (uint, error) {
			value, err := strconv.ParseUint(arg, 10, 0)
			return uint(value), err
		},
	}
}

func uintSliceParser() Parser[[]uint] {
	return &sliceArgParser[uint]{
		parser: func(arg string) (uint, error) {
			value, err := strconv.ParseUint(arg, 10, 0)
			return uint(value), err
		},
	}
}

func uint8Parser() Parser[uint8] {
	return &singleArgParser[uint8]{
		parser: func(arg string) (uint8, error) {
			value, err := strconv.ParseUint(arg, 10, 8)
			return uint8(value), err
		},
	}
}

func uint8SliceParser() Parser[[]uint8] {
	return &sliceArgParser[uint8]{
		parser: func(arg string) (uint8, error) {
			value, err := strconv.ParseUint(arg, 10, 8)
			return uint8(value), err
		},
	}
}

func uint16Parser() Parser[uint16] {
	return &singleArgParser[uint16]{
		parser: func(arg string) (uint16, error) {
			value, err := strconv.ParseUint(arg, 10, 16)
			return uint16(value), err
		},
	}
}

func uint16SliceParser() Parser[[]uint16] {
	return &sliceArgParser[uint16]{
		parser: func(arg string) (uint16, error) {
			value, err := strconv.ParseUint(arg, 10, 16)
			return uint16(value), err
		},
	}
}

func uint32Parser() Parser[uint32] {
	return &singleArgParser[uint32]{
		parser: func(arg string) (uint32, error) {
			value, err := strconv.ParseUint(arg, 10, 32)
			return uint32(value), err
		},
	}
}

func uint32SliceParser() Parser[[]uint32] {
	return &sliceArgParser[uint32]{
		parser: func(arg string) (uint32, error) {
			value, err := strconv.ParseUint(arg, 10, 32)
			return uint32(value), err
		},
	}
}

func uint64Parser() Parser[uint64] {
	return &singleArgParser[uint64]{
		parser: func(arg string) (uint64, error) {
			value, err := strconv.ParseUint(arg, 10, 64)
			return uint64(value), err
		},
	}
}

func uint64SliceParser() Parser[[]uint64] {
	return &sliceArgParser[uint64]{
		parser: func(arg string) (uint64, error) {
			value, err := strconv.ParseUint(arg, 10, 64)
			return uint64(value), err
		},
	}
}

func float32Parser() Parser[float32] {
	return &singleArgParser[float32]{
		parser: func(arg string) (float32, error) {
			value, err := strconv.ParseFloat(arg, 32)
			return float32(value), err
		},
	}
}

func float32SliceParser() Parser[[]float32] {
	return &sliceArgParser[float32]{
		parser: func(arg string) (float32, error) {
			value, err := strconv.ParseFloat(arg, 32)
			return float32(value), err
		},
	}
}

func float64Parser() Parser[float64] {
	return &singleArgParser[float64]{
		parser: func(arg string) (float64, error) {
			value, err := strconv.ParseFloat(arg, 64)
			return value, err
		},
	}
}

func float64SliceParser() Parser[[]float64] {
	return &sliceArgParser[float64]{
		parser: func(arg string) (float64, error) {
			value, err := strconv.ParseFloat(arg, 64)
			return value, err
		},
	}
}

func stringParser() Parser[string] {
	return &singleArgParser[string]{
		parser: func(arg string) (string, error) {
			return arg, nil
		},
	}
}

func stringSliceParser() Parser[[]string] {
	return &sliceArgParser[string]{
		parser: func(arg string) (string, error) {
			return arg, nil
		},
	}
}

type boolParser struct{}

func (p *boolParser) Parse(name string, args map[string][]string) (bool, error) {
	if len(args) > 1 {
		return false, errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %q", name)
	}

	for name, value := range args {
		if len(value) > 1 {
			return false, errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %q", name)
		}

		if strings.HasPrefix(name, "no-") {
			if len(value[0]) == 0 {
				return false, nil
			}

			value, err := strconv.ParseBool(value[0])
			if err != nil {
				return false, errors.Newf(err, ErrorCodeInvalidValue, "invalid value for flag %q", name)
			}
			return !value, nil
		}

		if len(value[0]) == 0 {
			return true, nil
		}
		value, err := strconv.ParseBool(value[0])
		if err != nil {
			return false, errors.Newf(err, ErrorCodeInvalidValue, "invalid value for flag %q", name)
		}
		return value, nil
	}
	return false, errors.Newf(nil, ErrorCodeMissingFlag, "missing flag %q", name)
}

type boolSliceParser struct{}

func (p *boolSliceParser) Parse(name string, args map[string][]string) ([]bool, error) {
	if len(args) == 0 {
		return nil, errors.Newf(nil, ErrorCodeMissingFlag, "missing flag %q", name)
	}

	var result []bool
	for name, value := range args {
		for _, value := range value {
			if strings.HasPrefix(name, "no-") {
				if len(value) == 0 {
					result = append(result, false)
					continue
				}

				value, err := strconv.ParseBool(value)
				if err != nil {
					return nil, errors.Newf(err, ErrorCodeInvalidValue, "invalid value for flag %q", name)
				}
				result = append(result, !value)
				continue
			}

			if len(value) == 0 {
				result = append(result, true)
				continue
			}
			value, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.Newf(err, ErrorCodeInvalidValue, "invalid value for flag %q", name)
			}
			result = append(result, value)
		}
	}
	return result, nil
}
