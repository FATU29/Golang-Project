package utils

import "os"

func StringValue(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func GetEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
