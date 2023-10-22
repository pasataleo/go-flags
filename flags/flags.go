package flags

import (
	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-inject/inject"
)

type Flags struct {
	flags   map[string]flag[any]
	aliases map[string]string
}

func (flags *Flags) Parse(args []string) error {

}

type FlagContext[T any] struct {
	flag flag[T]
}

func (ctx *FlagContext[T]) To(flags *Flags, target *T) error {
	if err := ctx.setFlag(flags); err != nil {
		return err
	}
	ctx.flag.Target = target
	return nil
}

func (ctx *FlagContext[T]) ToUnsafe(flags *Flags, target *T) {
	if err := ctx.To(flags, target); err != nil {
		panic(err)
	}
}

func (ctx *FlagContext[T]) ToInjector(flags *Flags, injector *inject.Injector, args ...string) error {
	if err := ctx.setFlag(flags); err != nil {
		return err
	}

	ctx.flag.Injector = injector
	ctx.flag.Args = args
	return nil
}

func (ctx *FlagContext[T]) ToInjectorUnsafe(flags *Flags, injector *inject.Injector, args ...string) {
	if err := ctx.ToInjector(flags, injector, args...); err != nil {
		panic(err)
	}
}

func (ctx *FlagContext[T]) setFlag(flags *Flags) error {
	if _, exists := flags.flags[ctx.flag.Name]; exists {
		return errors.Newf(nil, ErrorCodeDuplicateFlag, "Duplicate flag %s", ctx.flag.Name)
	}
	if _, exists := flags.aliases[ctx.flag.Name]; exists {
		return errors.Newf(nil, ErrorCodeDuplicateFlag, "Duplicate flag %s", ctx.flag.Name)
	}
	for _, alias := range ctx.flag.Aliases {
		if _, exists := flags.flags[alias]; exists {
			return errors.Newf(nil, ErrorCodeDuplicateFlag, "Duplicate flag %s", alias)
		}
		if _, exists := flags.aliases[alias]; exists {
			return errors.Newf(nil, ErrorCodeDuplicateFlag, "Duplicate flag %s", alias)
		}
	}

	flags.flags[ctx.flag.Name] = ctx.flag.(flag[interface{}])
	for _, alias := range flags.aliases {
		flags.aliases[alias] = ctx.flag.Name
	}
	return nil
}

type flag[T any] struct {
	Name     string
	Aliases  []string
	Parser   Parser[T]
	Default  T
	Optional bool

	// Injector and Args are used for injecting the flag value via an Injector.
	Injector *inject.Injector
	Args     []string

	// Target is used for injecting the flag value directly into a value.
	Target *T
}

func BindString(name string, optional bool, defaultValue string) FlagContext[string] {
	return FlagContext[string]{
		flag: flag[string]{
			Name:     name,
			Aliases:  []string{},
			Parser:   &StringParser{},
			Default:  defaultValue,
			Optional: optional,
		},
	}
}
