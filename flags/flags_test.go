package flags

import (
	"testing"

	"github.com/pasataleo/go-inject/inject"
	"github.com/pasataleo/go-testing/tests"
)

func TestFlags_Single(t *testing.T) {
	var value string

	flags := NewSet()

	BindString("value", "", false, "default").To(flags, &value)

	args := []string{"--value=hello"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, make([]string, 0))
	tests.Execute(value).Equal(t, "hello")
}

func TestFlags_Multi(t *testing.T) {
	var value []string

	flags := NewSet()

	BindStringSlice("value", "", false, nil).To(flags, &value)

	args := []string{"--value=hello", "--value=world"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, make([]string, 0))
	tests.Execute(value).Equal(t, []string{"hello", "world"})
}

func TestFlags_SingleAlternateFormat(t *testing.T) {
	var value string

	flags := NewSet()

	BindString("value", "", false, "default").To(flags, &value)

	args := []string{"--value", "hello"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, make([]string, 0))
	tests.Execute(value).Equal(t, "hello")
}

func TestFlags_SingleWithPath(t *testing.T) {
	var value string

	flags := NewSet()

	BindString("value", "", false, "default").To(flags, &value)

	args := []string{"--value", "hello", "world"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, []string{"world"})
	tests.Execute(value).Equal(t, "hello")
}

func TestFlags_SingleAlternateName(t *testing.T) {
	var value string

	flags := NewSet()

	BindString("value", "", false, "default").To(flags, &value)

	args := []string{"-value", "hello", "world"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, []string{"world"})
	tests.Execute(value).Equal(t, "hello")
}

func TestFlags_SingleOptional(t *testing.T) {
	var value string

	flags := NewSet()

	BindString("value", "", true, "default").To(flags, &value)

	args := []string{"world"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, []string{"world"})
	tests.Execute(value).Equal(t, "default")
}

func TestFlags_SingleAliased(t *testing.T) {
	var valueFalse bool
	var valueTrue bool

	flags := NewSet()

	BindBoolean("false", "", false, false).To(flags, &valueFalse)
	BindBoolean("true", "", false, true).To(flags, &valueTrue)

	args := []string{"--no-false", "--true"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, make([]string, 0))

	tests.Execute(valueFalse).Equal(t, false)
	tests.Execute(valueTrue).Equal(t, true)
}

func TestFlags_MultiAliased(t *testing.T) {
	var valueFalse []bool
	var valueTrue []bool

	flags := NewSet()

	BindBooleanSlice("false", "", false, nil).To(flags, &valueFalse)
	BindBooleanSlice("true", "", false, nil).To(flags, &valueTrue)

	args := []string{"--no-false", "--true", "--false"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, make([]string, 0))

	tests.Execute(valueFalse).Equal(t, []bool{false, true})
	tests.Execute(valueTrue).Equal(t, []bool{true})
}

func TestFlags_MissingRequiredFlag(t *testing.T) {
	var number int

	flags := NewSet()

	BindInt("number", "", false, 0).To(flags, &number)

	tests.Execute2E(flags.Parse([]string{})).ErrorCode(t, ErrorCodeMissingFlag)
}

func TestFlags_InvalidValue(t *testing.T) {
	var number int

	flags := NewSet()

	BindInt("number", "", false, 0).To(flags, &number)

	args := []string{"-number=notanumber"}
	tests.Execute2E(flags.Parse(args)).ErrorCode(t, ErrorCodeInvalidValue)
}

func TestFlags_Permissive(t *testing.T) {
	flags := NewSet()
	args := []string{"path/to/executable", "--value", "hello"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, args)
}

func TestFlags_ParseStrict(t *testing.T) {
	flags := NewSet()
	args := []string{"path/to/executable", "--value", "hello"}
	tests.Execute2E(flags.Parse(args, ParseBehaviorStrict)).ErrorCode(t, ErrorCodeUnknownFlag)
}

func TestFlags_ParseReadOnly(t *testing.T) {
	var value string

	flags := NewSet()

	BindString("value", "", false, "default").To(flags, &value)

	args := []string{"path/to/executable", "--value", "hello"}
	tests.Execute2E(flags.Parse(args, ParseBehaviorReadOnly)).NoError(t).Equal(t, args)
	tests.Execute(value).Equal(t, "hello")
}

func TestFlags_ParseToInjector(t *testing.T) {
	injector := inject.NewInjector()

	flags := NewSet()
	BindString("value", "", false, "default").ToInjector(flags, injector, "value")

	args := []string{"--value", "hello"}
	tests.Execute2E(flags.Parse(args)).NoError(t).Equal(t, make([]string, 0))
	tests.Execute2E(injector.Get("value")).NoError(t).Equal(t, "hello")
}
