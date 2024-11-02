package flags

import (
	"reflect"

	"github.com/pasataleo/go-errors/errors"
)

type Binder[T any] struct {
	flag *Flag[T]
}

type TargetFn[T any] func(name string, value T) error

func (binder *Binder[T]) ToValueSafe(flags *Set, target *T) error {
	binder.flag.target = reflect.ValueOf(target).Elem()
	if err := binder.setFlag(flags); err != nil {
		return err
	}
	return nil
}

func (binder *Binder[T]) ToValue(flags *Set, target *T) {
	if err := binder.ToValueSafe(flags, target); err != nil {
		panic(err)
	}
}

func (binder *Binder[T]) ToFunctionSafe(flags *Set, target TargetFn[T]) error {
	binder.flag.targetFn = target
	if err := binder.setFlag(flags); err != nil {
		return err
	}
	return nil
}

func (binder *Binder[T]) ToFunction(flags *Set, target TargetFn[T]) {
	if err := binder.ToFunctionSafe(flags, target); err != nil {
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
