package flags

import (
	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-inject/inject"
	"reflect"
)

type FlagContext[T any] struct {
	flag *flag[T]
}

func (ctx *FlagContext[T]) To(flags *Flags, target *T) error {
	ctx.flag.Target = reflect.ValueOf(target).Elem()
	if err := ctx.setFlag(flags); err != nil {
		return err
	}
	return nil
}

func (ctx *FlagContext[T]) ToUnsafe(flags *Flags, target *T) {
	if err := ctx.To(flags, target); err != nil {
		panic(err)
	}
}

func (ctx *FlagContext[T]) ToInjector(flags *Flags, injector *inject.Injector, args ...string) error {
	ctx.flag.Injector = injector
	ctx.flag.Args = args
	if err := ctx.setFlag(flags); err != nil {
		return err
	}
	return nil
}

func (ctx *FlagContext[T]) ToInjectorUnsafe(flags *Flags, injector *inject.Injector, args ...string) {
	if err := ctx.ToInjector(flags, injector, args...); err != nil {
		panic(err)
	}
}

func (ctx *FlagContext[T]) setFlag(flags *Flags) error {
	if _, exists := flags.flags[ctx.flag.Name]; exists {
		return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %s", ctx.flag.Name)
	}
	if _, exists := flags.aliases[ctx.flag.Name]; exists {
		return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %s", ctx.flag.Name)
	}
	for _, alias := range ctx.flag.Aliases {
		if _, exists := flags.flags[alias]; exists {
			return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %s", alias)
		}
		if _, exists := flags.aliases[alias]; exists {
			return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %s", alias)
		}
	}

	flags.flags[ctx.flag.Name] = ctx.flag.generic()
	for _, alias := range ctx.flag.Aliases {
		flags.aliases[alias] = ctx.flag.Name
	}
	return nil
}
