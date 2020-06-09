package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func GetEnvOrErr(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}

	return "", fmt.Errorf("ENV variable with key='%s' not set", key)
}

func GetEnvAsInt(key string, defaultVal int) int {
	valueStr := GetEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func GetEnvAsIntOrErr(key string) (int, error) {
	valueStr, err := GetEnvOrErr(key)
	if err != nil {
		return 0, fmt.Errorf("ENV variable with key='%s' not set", key)
	}

	if value, err := strconv.Atoi(valueStr); err != nil {
		return 0, fmt.Errorf("ENV variable with key='%s' can not be parsed to integer: %s", key, err)
	} else {
		return value, nil
	}
}
