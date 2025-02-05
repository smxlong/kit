package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_that_String_returns_the_value_when_set(t *testing.T) {
	key := "TEST_STRING"
	value := "test"
	defaultValue := "default"
	t.Setenv(key, value)
	require.Equal(t, value, String(key, defaultValue), "String() should return the value of the environment variable")
}

func Test_that_String_returns_the_default_value_when_not_set(t *testing.T) {
	key := "TEST_STRING"
	defaultValue := "default"
	require.Equal(t, defaultValue, String(key, defaultValue), "String() should return the default value")
}

func Test_that_RequiredString_returns_the_value_when_set(t *testing.T) {
	key := "TEST_REQUIRED_STRING"
	value := "test"
	t.Setenv(key, value)
	require.Equal(t, value, RequiredString(key), "RequiredString() should return the value of the environment variable")
}

func Test_that_RequiredString_panics_when_not_set(t *testing.T) {
	key := "TEST_REQUIRED_STRING"
	require.Panics(t, func() { RequiredString(key) }, "RequiredString() should panic when the environment variable is not set")
}

func Test_that_Int_returns_the_value_when_set(t *testing.T) {
	key := "TEST_INT"
	value := "42"
	defaultValue := 0
	t.Setenv(key, value)
	require.Equal(t, 42, Int(key, defaultValue), "Int() should return the value of the environment variable")
}

func Test_that_Int_returns_the_default_value_when_not_set(t *testing.T) {
	key := "TEST_INT"
	defaultValue := 42
	require.Equal(t, defaultValue, Int(key, defaultValue), "Int() should return the default value")
}

func Test_that_Int_returns_the_default_value_when_not_a_valid_integer(t *testing.T) {
	key := "TEST_INT"
	value := "not an integer"
	defaultValue := 42
	t.Setenv(key, value)
	require.Equal(t, defaultValue, Int(key, defaultValue), "Int() should return the default value")
}

func Test_that_RequiredInt_returns_the_value_when_set(t *testing.T) {
	key := "TEST_REQUIRED_INT"
	value := "42"
	t.Setenv(key, value)
	require.Equal(t, 42, RequiredInt(key), "RequiredInt() should return the value of the environment variable")
}

func Test_that_RequiredInt_panics_when_not_set(t *testing.T) {
	key := "TEST_REQUIRED_INT"
	require.Panics(t, func() { RequiredInt(key) }, "RequiredInt() should panic when the environment variable is not set")
}

func Test_that_RequiredInt_panics_when_not_a_valid_integer(t *testing.T) {
	key := "TEST_REQUIRED_INT"
	value := "not an integer"
	t.Setenv(key, value)
	require.Panics(t, func() { RequiredInt(key) }, "RequiredInt() should panic when the environment variable is not a valid integer")
}

func Test_that_Bool_returns_the_value_when_set(t *testing.T) {
	key := "TEST_BOOL"
	value := "true"
	defaultValue := false
	t.Setenv(key, value)
	require.True(t, Bool(key, defaultValue), "Bool() should return the value of the environment variable")
}

func Test_that_Bool_returns_the_default_value_when_not_set(t *testing.T) {
	key := "TEST_BOOL"
	defaultValue := true
	require.True(t, Bool(key, defaultValue), "Bool() should return the default value")
}

func Test_that_Bool_returns_the_default_value_when_not_a_valid_boolean(t *testing.T) {
	key := "TEST_BOOL"
	value := "not a boolean"
	defaultValue := true
	t.Setenv(key, value)
	require.True(t, Bool(key, defaultValue), "Bool() should return the default value")
}

func Test_that_RequiredBool_returns_the_value_when_set(t *testing.T) {
	key := "TEST_REQUIRED_BOOL"
	value := "true"
	t.Setenv(key, value)
	require.True(t, RequiredBool(key), "RequiredBool() should return the value of the environment variable")
}

func Test_that_RequiredBool_panics_when_not_set(t *testing.T) {
	key := "TEST_REQUIRED_BOOL"
	require.Panics(t, func() { RequiredBool(key) }, "RequiredBool() should panic when the environment variable is not set")
}

func Test_that_RequiredBool_panics_when_not_a_valid_boolean(t *testing.T) {
	key := "TEST_REQUIRED_BOOL"
	value := "not a boolean"
	t.Setenv(key, value)
	require.Panics(t, func() { RequiredBool(key) }, "RequiredBool() should panic when the environment variable is not a valid boolean")
}

func Test_that_Float64_returns_the_value_when_set(t *testing.T) {
	key := "TEST_FLOAT64"
	value := "42.42"
	defaultValue := 0.0
	t.Setenv(key, value)
	require.Equal(t, 42.42, Float64(key, defaultValue), "Float64() should return the value of the environment variable")
}

func Test_that_Float64_returns_the_default_value_when_not_set(t *testing.T) {
	key := "TEST_FLOAT64"
	defaultValue := 42.42
	require.Equal(t, defaultValue, Float64(key, defaultValue), "Float64() should return the default value")
}

func Test_that_Float64_returns_the_default_value_when_not_a_valid_float(t *testing.T) {
	key := "TEST_FLOAT64"
	value := "not a float"
	defaultValue := 42.42
	t.Setenv(key, value)
	require.Equal(t, defaultValue, Float64(key, defaultValue), "Float64() should return the default value")
}

func Test_that_RequiredFloat64_returns_the_value_when_set(t *testing.T) {
	key := "TEST_REQUIRED_FLOAT64"
	value := "42.42"
	t.Setenv(key, value)
	require.Equal(t, 42.42, RequiredFloat64(key), "RequiredFloat64() should return the value of the environment variable")
}

func Test_that_RequiredFloat64_panics_when_not_set(t *testing.T) {
	key := "TEST_REQUIRED_FLOAT64"
	require.Panics(t, func() { RequiredFloat64(key) }, "RequiredFloat64() should panic when the environment variable is not set")
}

func Test_that_RequiredFloat64_panics_when_not_a_valid_float(t *testing.T) {
	key := "TEST_REQUIRED_FLOAT64"
	value := "not a float"
	t.Setenv(key, value)
	require.Panics(t, func() { RequiredFloat64(key) }, "RequiredFloat64() should panic when the environment variable is not a valid float")
}

func Test_that_Strings_returns_the_value_when_set(t *testing.T) {
	key := "TEST_STRINGS"
	value := "test1,test2,test3"
	defaultValue := []string{"default"}
	t.Setenv(key, value)
	require.Equal(t, []string{"test1", "test2", "test3"}, Strings(key, ",", defaultValue), "Strings() should return the value of the environment variable")
}

func Test_that_Strings_returns_the_default_value_when_not_set(t *testing.T) {
	key := "TEST_STRINGS"
	defaultValue := []string{"default"}
	require.Equal(t, defaultValue, Strings(key, ",", defaultValue), "Strings() should return the default value")
}

func Test_that_RequiredStrings_returns_the_value_when_set(t *testing.T) {
	key := "TEST_REQUIRED_STRINGS"
	value := "test1,test2,test3"
	t.Setenv(key, value)
	require.Equal(t, []string{"test1", "test2", "test3"}, RequiredStrings(key, ","), "RequiredStrings() should return the value of the environment variable")
}

func Test_that_RequiredStrings_panics_when_not_set(t *testing.T) {
	key := "TEST_REQUIRED_STRINGS"
	require.Panics(t, func() { RequiredStrings(key, ",") }, "RequiredStrings() should panic when the environment variable is not set")
}

func Test_that_RequiredNonEmptyStrings_returns_the_value_when_set(t *testing.T) {
	key := "TEST_REQUIRED_NON_EMPTY_STRINGS"
	value := "test1,test2,test3"
	t.Setenv(key, value)
	require.Equal(t, []string{"test1", "test2", "test3"}, RequiredNonEmptyStrings(key, ","), "RequiredNonEmptyStrings() should return the value of the environment variable")
}

func Test_that_RequiredNonEmptyStrings_panics_when_not_set(t *testing.T) {
	key := "TEST_REQUIRED_NON_EMPTY_STRINGS"
	require.Panics(t, func() { RequiredNonEmptyStrings(key, ",") }, "RequiredNonEmptyStrings() should panic when the environment variable is not set")
}

func Test_that_RequiredNonEmptyStrings_panics_when_empty(t *testing.T) {
	key := "TEST_REQUIRED_NON_EMPTY_STRINGS"
	value := ""
	t.Setenv(key, value)
	require.Panics(t, func() { RequiredNonEmptyStrings(key, ",") }, "RequiredNonEmptyStrings() should panic when the environment variable is an empty string")
}

func Test_that_RequiredNonEmptyStrings_panics_when_split_values_are_empty(t *testing.T) {
	key := "TEST_REQUIRED_NON_EMPTY_STRINGS"
	value := "test1,,test3"
	t.Setenv(key, value)
	require.Panics(t, func() { RequiredNonEmptyStrings(key, ",") }, "RequiredNonEmptyStrings() should panic when the split values are empty strings")
}
