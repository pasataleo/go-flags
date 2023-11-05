package flags

import (
	"testing"

	"github.com/pasataleo/go-inject/inject"
	"github.com/pasataleo/go-testing/tests"
)

func TestFlags_Single(t *testing.T) {
	var value string

	flags := NewSet()

	ctx := BindString("value", "", false, "default")
	ctx.ToUnsafe(flags, &value)

	args := []string{"--value=hello"}
	tests.ExecFn(t, flags.Parse, args).NoError()
	tests.Value(t, value).Equals("hello")
}

func TestFlags_SingleAlternateFormat(t *testing.T) {
	var value string

	flags := NewSet()

	ctx := BindString("value", "", false, "default")
	ctx.ToUnsafe(flags, &value)

	args := []string{"--value", "hello"}
	tests.ExecFn(t, flags.Parse, args).NoError()
	tests.Value(t, value).Equals("hello")
}

func TestFlags_SingleWithPath(t *testing.T) {
	var value string

	flags := NewSet()

	ctx := BindString("value", "", false, "default")
	ctx.ToUnsafe(flags, &value)

	args := []string{"--value", "hello", "world"}
	tests.ExecFn(t, flags.Parse, args).NoError().Equals([]string{"world"})
	tests.Value(t, value).Equals("hello")
}

func TestFlags_SingleAlternateName(t *testing.T) {
	var value string

	flags := NewSet()

	ctx := BindString("value", "", false, "default")
	ctx.ToUnsafe(flags, &value)

	args := []string{"-value", "hello", "world"}
	tests.ExecFn(t, flags.Parse, args).NoError().Equals([]string{"world"})
	tests.Value(t, value).Equals("hello")
}

func TestFlags_SingleOptional(t *testing.T) {
	var value string

	flags := NewSet()

	ctx := BindString("value", "", true, "default")
	ctx.ToUnsafe(flags, &value)

	tests.ExecFn(t, flags.Parse, nil).NoError().Empty()
	tests.Value(t, value).Equals("default")
}

func TestFlags_SingleAliased(t *testing.T) {
	var valueFalse bool
	var valueTrue bool

	flags := NewSet()

	BindBoolean("false", "", false, false).ToUnsafe(flags, &valueFalse)
	BindBoolean("true", "", false, true).ToUnsafe(flags, &valueTrue)

	args := []string{"--no-false", "--true"}
	tests.ExecFn(t, flags.Parse, args).NoError().Empty()

	tests.Value(t, valueFalse).Equals(false)
	tests.Value(t, valueTrue).Equals(true)
}

func TestFlags_MissingRequiredFlag(t *testing.T) {
	var number int

	flags := NewSet()

	BindInt("number", "", false, 0).ToUnsafe(flags, &number)

	tests.ExecFn(t, flags.Parse, nil).ErrorCode(ErrorCodeMissingFlag)
}

func TestFlags_InvalidValue(t *testing.T) {
	var number int

	flags := NewSet()

	BindInt("number", "", false, 0).ToUnsafe(flags, &number)

	args := []string{"-number=notanumber"}
	tests.ExecFn(t, flags.Parse, args).ErrorCode(ErrorCodeInvalidValue)
}

func TestFlags_Permissive(t *testing.T) {
	flags := NewSet()
	args := []string{"path/to/executable", "--value", "hello"}
	tests.ExecFn(t, flags.Parse, args).NoError().Equals([]string{"path/to/executable", "--value", "hello"})
}

func TestFlags_ParseStrict(t *testing.T) {
	flags := NewSet()
	args := []string{"path/to/executable", "--value", "hello"}
	tests.ExecFn(t, flags.Parse, args, ParseBehaviorStrict).ErrorCode(ErrorCodeUnknownFlag)
}

func TestFlags_ParseReadOnly(t *testing.T) {
	var value string

	flags := NewSet()

	BindString("value", "", false, "default").ToUnsafe(flags, &value)

	args := []string{"path/to/executable", "--value", "hello"}
	tests.ExecFn(t, flags.Parse, args, ParseBehaviorReadOnly).NoError().Equals(args)
}

func TestFlags_ParseToInjector(t *testing.T) {
	injector := inject.NewInjector()

	flags := NewSet()
	BindString("value", "", false, "default").ToInjectorUnsafe(flags, injector, "value")

	args := []string{"--value", "hello"}
	tests.ExecFn(t, flags.Parse, args).NoError()

	tests.ExecFn(t, injector.Get, "value").NoError().Equals("hello")
}
