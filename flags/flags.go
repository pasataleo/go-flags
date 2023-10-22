package flags

import (
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

	flags.flags[ctx.flag.Aliases[0]] = ctx.flag
}

func (ctx *FlagContext[T]) ToUnsafe(target *T) {

}

func (ctx *FlagContext[T]) ToInjector(injector *inject.Injector) error {

}

func (ctx *FlagContext[T]) ToInjectorUnsafe(injector *inject.Injector) {

}

type flag[T any] struct {
	Name     string
	Aliases  []string
	Parser   Parser[T]
	Default  T
	Optional bool

	Injector *inject.Injector
	Target   *T
}

func BindString(name string, optional bool, defaultValue string) FlagContext[string] {
	return FlagContext[string]{
		flag: flag[string]{
			Aliases:  []string{name},
			Parser:   &StringParser{},
			Default:  defaultValue,
			Optional: optional,
		},
	}
}
