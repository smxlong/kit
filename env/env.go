package env

import (
	"os"
	"strconv"
	"strings"
)

// String returns the value of the environment variable named by the key, or
// the default value if the environment variable is not set.
func String(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// RequiredString returns the value of the environment variable named by the key,
// or panics if the environment variable is not set.
func RequiredString(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("required environment variable " + key + " is not set")
}

// Int returns the value of the environment variable named by the key, or the
// default value if the environment variable is not set or is not a valid integer.
func Int(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// RequiredInt returns the value of the environment variable named by the key,
// or panics if the environment variable is not set or is not a valid integer.
func RequiredInt(key string) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		panic("required environment variable " + key + " is not a valid integer")
	}
	panic("required environment variable " + key + " is not set")
}

// Bool returns the value of the environment variable named by the key, or the
// default value if the environment variable is not set or is not a valid boolean.
func Bool(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// RequiredBool returns the value of the environment variable named by the key,
// or panics if the environment variable is not set or is not a valid boolean.
func RequiredBool(key string) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		panic("required environment variable " + key + " is not a valid boolean")
	}
	panic("required environment variable " + key + " is not set")
}

// Float64 returns the value of the environment variable named by the key, or the
// default value if the environment variable is not set or is not a valid float64.
func Float64(key string, defaultValue float64) float64 {
	if value, ok := os.LookupEnv(key); ok {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// RequiredFloat64 returns the value of the environment variable named by the key,
// or panics if the environment variable is not set or is not a valid float64.
func RequiredFloat64(key string) float64 {
	if value, ok := os.LookupEnv(key); ok {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
		panic("required environment variable " + key + " is not a valid float64")
	}
	panic("required environment variable " + key + " is not set")
}

// Strings returns the value of the environment variable named by the key, or the
// default value if the environment variable is not set. The value is split by
// the specified separator.
func Strings(key, separator string, defaultValue []string) []string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.Split(value, separator)
	}
	return defaultValue
}

// RequiredStrings returns the value of the environment variable named by the key,
// or panics if the environment variable is not set. The value is split by the
// specified separator.
func RequiredStrings(key, separator string) []string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.Split(value, separator)
	}
	panic("required environment variable " + key + " is not set")
}

// RequiredNonEmptyStrings returns the value of the environment variable named by
// the key, or panics if the environment variable is not set or is an empty string,
// or if any of the split values are empty strings. The value is split by the
// specified separator.
func RequiredNonEmptyStrings(key, separator string) []string {
	if value, ok := os.LookupEnv(key); ok {
		if value == "" {
			panic("required environment variable " + key + " is an empty string")
		}
		values := strings.Split(value, separator)
		for _, v := range values {
			if v == "" {
				panic("required environment variable " + key + " contains an empty string")
			}
		}
		return values
	}
	panic("required environment variable " + key + " is not set")
}
