package environment

import (
	"os"
	"testing"
)

// GetFirstNonEmptyEnvVarOrFatal returns the first non-empty environment variable from envVarNames, or throws a fatal
func GetFirstNonEmptyEnvVarOrFatal(t *testing.T, envVarNames []string) string {
	value := GetFirstNonEmptyEnvVarOrEmptyString(t, envVarNames)
	if value == "" {
		t.Fatalf("All of the following env vars %v are empty. At least one must be non-empty.", envVarNames)
	}

	return value
}

// GetFirstNonEmptyEnvVarOrEmptyString returns the first non-empty environment variable from envVarNames, or returns the
// empty string
func GetFirstNonEmptyEnvVarOrEmptyString(t *testing.T, envVarNames []string) string {
	for _, name := range envVarNames {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}

	return ""
}
