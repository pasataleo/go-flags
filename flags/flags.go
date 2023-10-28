package flags

import (
	"fmt"
	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-inject/inject"
	"reflect"
	"strings"
)

type ParseBehavior int

const (
	ParseBehaviorStrict = iota
	ParseBehaviorReadOnly
)

type Flags struct {
	flags   map[string]*flag[any]
	aliases map[string]string
}

func Set() *Flags {
	return &Flags{
		flags:   make(map[string]*flag[any]),
		aliases: make(map[string]string),
	}
}

func (flags *Flags) Parse(args []string, behaviours ...ParseBehavior) ([]string, error) {
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

func (flags *Flags) parse(args []string, strict bool, readOnly bool) ([]string, error) {
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

		if _, exists := flags.flags[name]; !exists {
			if strict {
				err = errors.Append(err, errors.Newf(nil, ErrorCodeUnknownFlag, "unknown flag %s", name))
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

	for name, flag := range flags.flags {
		if _, exists := unparsed[name]; !exists {
			if !flag.Optional {
				err = errors.Append(err, errors.Newf(nil, ErrorCodeMissingFlag, "missing flag %s", name))
				continue
			}

			if err := flag.setValue(flag.Default); err != nil {
				err = errors.Append(err, errors.Newf(err, ErrorCodeInvalidValue, "could not set default value for %s", name))
			}
			continue
		}

		values := unparsed[name]
		delete(unparsed, name)

		if flag.Parser != nil {
			var flattened []string
			for _, value := range values {
				flattened = append(flattened, value)
			}

			value, valueErr := flag.Parser.Parse(name, flattened)
			if valueErr != nil {
				err = errors.Append(err, errors.Newf(valueErr, ErrorCodeInvalidValue, "invalid flag %s", name))
				continue
			}

			if err := flag.setValue(value); err != nil {
				err = errors.Append(err, errors.Newf(err, ErrorCodeInvalidValue, "invalid flag %s", name))
			}

			continue
		}

		if flag.AliasParser != nil {
			value, valueErr := flag.AliasParser.Parse(name, values)
			if valueErr != nil {
				err = errors.Append(err, errors.Newf(valueErr, ErrorCodeInvalidValue, "invalid flag %s", name))
				continue
			}

			if err := flag.setValue(value); err != nil {
				err = errors.Append(err, errors.Newf(err, ErrorCodeInvalidValue, "invalid flag %s", name))
			}

			continue
		}

		// We should ensure all flags have a parser or alias parser.
		panic("flag doesn't have a parser")
	}

	return remaining, err
}

type flag[T any] struct {
	Name     string
	Aliases  []string
	Default  T
	Optional bool

	Parser      Parser[T]
	AliasParser aliasParser[T]

	// Injector and Args are used for injecting the flag value via an Injector.
	Injector *inject.Injector
	Args     []string

	// Target is used for injecting the flag value directly into a value.
	Target reflect.Value
}

func (f *flag[T]) setValue(value T) error {
	if f.Injector != nil {
		return inject.BindValue(value).To(f.Injector, f.Args...)
	}

	f.Target.Set(reflect.ValueOf(value))
	return nil
}

func (f *flag[T]) generic() *flag[interface{}] {
	generic := &flag[interface{}]{
		Name:     f.Name,
		Aliases:  f.Aliases,
		Default:  f.Default,
		Optional: f.Optional,
		Injector: f.Injector,
		Args:     f.Args,
		Target:   f.Target,
	}

	if f.Parser != nil {
		generic.Parser = &parserWrapper[T]{
			parser: f.Parser,
		}
	}

	if f.AliasParser != nil {
		generic.AliasParser = &aliasWrapper[T]{
			parser: f.AliasParser,
		}
	}

	return generic
}

func BindValue[T any](name string, optional bool, defaultValue T, Parser Parser[T]) *FlagContext[T] {
	return &FlagContext[T]{
		flag: &flag[T]{
			Name:     name,
			Parser:   Parser,
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindBoolean(name string, optional bool, defaultValue bool) *FlagContext[bool] {
	return &FlagContext[bool]{
		flag: &flag[bool]{
			Name: name,
			Aliases: []string{
				fmt.Sprintf("no-%s", name),
			},
			AliasParser: &boolParser{},
			Default:     defaultValue,
			Optional:    optional,
		},
	}
}

func BindString(name string, optional bool, defaultValue string) *FlagContext[string] {
	return &FlagContext[string]{
		flag: &flag[string]{
			Name:     name,
			Parser:   stringParser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindInt(name string, optional bool, defaultValue int) *FlagContext[int] {
	return &FlagContext[int]{
		flag: &flag[int]{
			Name:     name,
			Parser:   intParser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindInt8(name string, optional bool, defaultValue int8) *FlagContext[int8] {
	return &FlagContext[int8]{
		flag: &flag[int8]{
			Name:     name,
			Parser:   int8Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindInt16(name string, optional bool, defaultValue int16) *FlagContext[int16] {
	return &FlagContext[int16]{
		flag: &flag[int16]{
			Name:     name,
			Parser:   int16Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindInt32(name string, optional bool, defaultValue int32) *FlagContext[int32] {
	return &FlagContext[int32]{
		flag: &flag[int32]{
			Name:     name,
			Parser:   int32Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindInt64(name string, optional bool, defaultValue int64) *FlagContext[int64] {
	return &FlagContext[int64]{
		flag: &flag[int64]{
			Name:     name,
			Parser:   int64Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindUint(name string, optional bool, defaultValue uint) *FlagContext[uint] {
	return &FlagContext[uint]{
		flag: &flag[uint]{
			Name:     name,
			Parser:   uintParser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindUint8(name string, optional bool, defaultValue uint8) *FlagContext[uint8] {
	return &FlagContext[uint8]{
		flag: &flag[uint8]{
			Name:     name,
			Parser:   uint8Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindUint16(name string, optional bool, defaultValue uint16) *FlagContext[uint16] {
	return &FlagContext[uint16]{
		flag: &flag[uint16]{
			Name:     name,
			Parser:   uint16Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindUint32(name string, optional bool, defaultValue uint32) *FlagContext[uint32] {
	return &FlagContext[uint32]{
		flag: &flag[uint32]{
			Name:     name,
			Parser:   uint32Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindUint64(name string, optional bool, defaultValue uint64) *FlagContext[uint64] {
	return &FlagContext[uint64]{
		flag: &flag[uint64]{
			Name:     name,
			Parser:   uint64Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindFloat32(name string, optional bool, defaultValue float32) *FlagContext[float32] {
	return &FlagContext[float32]{
		flag: &flag[float32]{
			Name:     name,
			Parser:   float32Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}

func BindFloat64(name string, optional bool, defaultValue float64) *FlagContext[float64] {
	return &FlagContext[float64]{
		flag: &flag[float64]{
			Name:     name,
			Parser:   float64Parser(),
			Default:  defaultValue,
			Optional: optional,
		},
	}
}
