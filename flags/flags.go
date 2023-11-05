package flags

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-inject/inject"
)

type ParseBehavior int

const (
	ParseBehaviorStrict = iota
	ParseBehaviorReadOnly
)

type Set struct {
	Flags   map[string]*Flag[any]
	aliases map[string]string
}

func NewSet() *Set {
	return &Set{
		Flags:   make(map[string]*Flag[any]),
		aliases: make(map[string]string),
	}
}

func (flags *Set) Parse(args []string, behaviours ...ParseBehavior) ([]string, error) {
	var strict, readOnly bool
	for _, behaviour := range behaviours {
		switch behaviour {
		case ParseBehaviorStrict:
			strict = true
		case ParseBehaviorReadOnly:
			readOnly = true
		}
	}

	return flags.parse(args, strict, readOnly)
}

func (flags *Set) parse(args []string, strict bool, readOnly bool) ([]string, error) {
	var err error

	// Anything we don't process will be returned.
	var remaining []string
	if readOnly {
		remaining = args
	}

	appendRemaining := func(arg string) {
		if readOnly {
			return
		}
		remaining = append(remaining, arg)
	}

	// unparsed maps flag names to unparsed flag values. We can have multiple values for a flag and aliases for flag
	// names.
	unparsed := make(map[string]map[string]string)

	isFlagName := func(arg string) (string, bool) {
		if name, ok := strings.CutPrefix(arg, "--"); ok {
			return name, true
		}

		if name, ok := strings.CutPrefix(arg, "-"); ok {
			return name, true
		}

		return arg, false
	}

	hasFlagValue := func(arg string) (string, string, bool) {
		if name, value, ok := strings.Cut(arg, "="); ok {
			return name, value, true
		}
		return arg, "", false
	}

	resolveName := func(name string) (string, string) {
		if alias, exists := flags.aliases[name]; exists {
			return alias, name
		}
		return name, name
	}

	skipNextArg := false
	for ix, arg := range args {
		if skipNextArg {
			skipNextArg = false
			continue
		}

		flag, isFlag := isFlagName(arg)
		if !isFlag {
			appendRemaining(arg)
			continue
		}

		name, value, containsValue := hasFlagValue(flag)
		name, alias := resolveName(name)

		if _, exists := flags.Flags[name]; !exists {
			if strict {
				err = errors.Append(err, errors.Newf(nil, ErrorCodeUnknownFlag, "unknown flag %q", name))
			}
			appendRemaining(arg)
			continue
		}

		if containsValue {
			if _, exists := unparsed[name]; !exists {
				unparsed[name] = make(map[string]string)
			}
			unparsed[name][alias] = value
			continue
		}

		if ix+1 < len(args) {
			nextArg := args[ix+1]
			if _, isFlag := isFlagName(nextArg); !isFlag {
				if _, exists := unparsed[name]; !exists {
					unparsed[name] = make(map[string]string)
				}
				unparsed[name][alias] = nextArg
				skipNextArg = true
				continue
			}
		}

		if _, exists := unparsed[name]; !exists {
			unparsed[name] = make(map[string]string)
		}
		unparsed[name][alias] = ""
	}

	for name, flag := range flags.Flags {
		if _, exists := unparsed[name]; !exists {
			if !flag.Optional {
				err = errors.Append(err, errors.Newf(nil, ErrorCodeMissingFlag, "missing flag %q", name))
				continue
			}

			if valueErr := flag.setValue(flag.Default); valueErr != nil {
				err = errors.Append(err, errors.Newf(valueErr, ErrorCodeInvalidValue, "could not set default value for %q", name))
			}
			continue
		}

		values := unparsed[name]
		delete(unparsed, name)

		if flag.parser != nil {
			var flattened []string
			for _, value := range values {
				flattened = append(flattened, value)
			}

			value, valueErr := flag.parser.Parse(name, flattened)
			if valueErr != nil {
				err = errors.Append(err, errors.Newf(valueErr, ErrorCodeInvalidValue, "invalid flag %q", name))
				continue
			}

			if valueErr := flag.setValue(value); valueErr != nil {
				err = errors.Append(err, errors.Newf(valueErr, ErrorCodeInvalidValue, "invalid flag %q", name))
			}

			continue
		}

		if flag.aliasParser != nil {
			value, valueErr := flag.aliasParser.Parse(name, values)
			if valueErr != nil {
				err = errors.Append(err, errors.Newf(valueErr, ErrorCodeInvalidValue, "invalid flag %q", name))
				continue
			}

			if valueErr := flag.setValue(value); valueErr != nil {
				err = errors.Append(err, errors.Newf(valueErr, ErrorCodeInvalidValue, "invalid flag %q", name))
			}

			continue
		}

		// We should ensure all flags have a parser or alias parser.
		panic("flag doesn't have a parser")
	}

	return remaining, err
}

type Flag[T any] struct {
	Name        string
	Aliases     []string
	Default     T
	Optional    bool
	Description string

	parser      Parser[T]
	aliasParser aliasParser[T]

	// injector and args are used for injecting the flag value via an injector.
	injector *inject.Injector
	args     []string

	// target is used for injecting the flag value directly into a value.
	target reflect.Value
}

func (f *Flag[T]) setValue(value T) error {
	if f.injector != nil {
		return inject.BindValue(value).To(f.injector, f.args...)
	}

	f.target.Set(reflect.ValueOf(value))
	return nil
}

func (f *Flag[T]) generic() *Flag[interface{}] {
	generic := &Flag[interface{}]{
		Name:        f.Name,
		Aliases:     f.Aliases,
		Default:     f.Default,
		Optional:    f.Optional,
		Description: f.Description,
		injector:    f.injector,
		args:        f.args,
		target:      f.target,
	}

	if f.parser != nil {
		generic.parser = &parserWrapper[T]{
			parser: f.parser,
		}
	}

	if f.aliasParser != nil {
		generic.aliasParser = &aliasWrapper[T]{
			parser: f.aliasParser,
		}
	}

	return generic
}

func BindValue[T any](name string, description string, optional bool, defaultValue T, parser Parser[T]) *Binder[T] {
	return &Binder[T]{
		flag: &Flag[T]{
			Name:        name,
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
			parser:      parser,
		},
	}
}

func BindBoolean(name string, description string, optional bool, defaultValue bool) *Binder[bool] {
	return &Binder[bool]{
		flag: &Flag[bool]{
			Name: name,
			Aliases: []string{
				fmt.Sprintf("no-%s", name),
			},
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
			aliasParser: &boolParser{},
		},
	}
}

func BindString(name string, description string, optional bool, defaultValue string) *Binder[string] {
	return &Binder[string]{
		flag: &Flag[string]{
			Name:        name,
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
			parser:      stringParser(),
		},
	}
}

func BindInt(name string, description string, optional bool, defaultValue int) *Binder[int] {
	return &Binder[int]{
		flag: &Flag[int]{
			Name:        name,
			parser:      intParser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindInt8(name string, description string, optional bool, defaultValue int8) *Binder[int8] {
	return &Binder[int8]{
		flag: &Flag[int8]{
			Name:        name,
			parser:      int8Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindInt16(name string, description string, optional bool, defaultValue int16) *Binder[int16] {
	return &Binder[int16]{
		flag: &Flag[int16]{
			Name:        name,
			parser:      int16Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindInt32(name string, description string, optional bool, defaultValue int32) *Binder[int32] {
	return &Binder[int32]{
		flag: &Flag[int32]{
			Name:        name,
			parser:      int32Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindInt64(name string, description string, optional bool, defaultValue int64) *Binder[int64] {
	return &Binder[int64]{
		flag: &Flag[int64]{
			Name:        name,
			parser:      int64Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindUint(name string, description string, optional bool, defaultValue uint) *Binder[uint] {
	return &Binder[uint]{
		flag: &Flag[uint]{
			Name:        name,
			parser:      uintParser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindUint8(name string, description string, optional bool, defaultValue uint8) *Binder[uint8] {
	return &Binder[uint8]{
		flag: &Flag[uint8]{
			Name:        name,
			parser:      uint8Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindUint16(name string, description string, optional bool, defaultValue uint16) *Binder[uint16] {
	return &Binder[uint16]{
		flag: &Flag[uint16]{
			Name:        name,
			parser:      uint16Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindUint32(name string, description string, optional bool, defaultValue uint32) *Binder[uint32] {
	return &Binder[uint32]{
		flag: &Flag[uint32]{
			Name:        name,
			parser:      uint32Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindUint64(name string, description string, optional bool, defaultValue uint64) *Binder[uint64] {
	return &Binder[uint64]{
		flag: &Flag[uint64]{
			Name:        name,
			parser:      uint64Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindFloat32(name string, description string, optional bool, defaultValue float32) *Binder[float32] {
	return &Binder[float32]{
		flag: &Flag[float32]{
			Name:        name,
			parser:      float32Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}

func BindFloat64(name string, description string, optional bool, defaultValue float64) *Binder[float64] {
	return &Binder[float64]{
		flag: &Flag[float64]{
			Name:        name,
			parser:      float64Parser(),
			Default:     defaultValue,
			Optional:    optional,
			Description: description,
		},
	}
}
