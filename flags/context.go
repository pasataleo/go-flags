package flags

import (
	"reflect"

	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-inject/inject"
)

type Binder[T any] struct {
	flag *Flag[T]
}

func (binder *Binder[T]) ToSafe(flags *Set, target *T) error {
	binder.flag.target = reflect.ValueOf(target).Elem()
	if err := binder.setFlag(flags); err != nil {
		return err
	}
	return nil
}

func (binder *Binder[T]) To(flags *Set, target *T) {
	if err := binder.ToSafe(flags, target); err != nil {
		panic(err)
	}
}

func (binder *Binder[T]) ToInjectorSafe(flags *Set, injector *inject.Injector, args ...string) error {
	binder.flag.injector = injector
	binder.flag.args = args
	if err := binder.setFlag(flags); err != nil {
		return err
	}
	return nil
}

func (binder *Binder[T]) ToInjector(flags *Set, injector *inject.Injector, args ...string) {
	if err := binder.ToInjectorSafe(flags, injector, args...); err != nil {
		panic(err)
	}
}

func (binder *Binder[T]) setFlag(flags *Set) error {
	if _, exists := flags.Flags[binder.flag.Name]; exists {
		return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %q", binder.flag.Name)
	}
	if _, exists := flags.aliases[binder.flag.Name]; exists {
		return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %q", binder.flag.Name)
	}
	for _, alias := range binder.flag.Aliases {
		if _, exists := flags.Flags[alias]; exists {
			return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %q", alias)
		}
		if _, exists := flags.aliases[alias]; exists {
			return errors.Newf(nil, ErrorCodeDuplicateFlag, "duplicate flag %q", alias)
		}
	}

	flags.Flags[binder.flag.Name] = binder.flag.generic()
	for _, alias := range binder.flag.Aliases {
		flags.aliases[alias] = binder.flag.Name
	}
	return nil
}
