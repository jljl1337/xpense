package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MustGetBool(key string, defaultValue bool) bool {
	value, err := GetBool(key, defaultValue)
	if err != nil {
		panic(err)
	}
	return value
}

func GetBool(key string, defaultValue bool) (bool, error) {
	defaultStr := "false"
	if defaultValue {
		defaultStr = "true"
	}

	value, err := GetString(key, defaultStr)
	if err != nil {
		return false, err
	}

	return strings.ToLower(value) == "true", nil
}

func MustGetInt(key string, defaultValue int) int {
	value, err := GetInt(key, defaultValue)
	if err != nil {
		panic(err)
	}
	return value
}

func GetInt(key string, defaultValue int) (int, error) {
	value, err := GetString(key, fmt.Sprintf("%d", defaultValue))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func MustGetString(key string, defaultValue string) string {
	value, err := GetString(key, defaultValue)
	if err != nil {
		panic(err)
	}
	return value
}

func GetString(key string, defaultValue string) (string, error) {
	fileKey := key + "_FILE"
	fileValue, fileExists := os.LookupEnv(fileKey)
	if fileExists {
		// Read from the file
		data, err := os.ReadFile(fileValue)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	value, exists := os.LookupEnv(key)
	if exists {
		return value, nil
	}
	return defaultValue, nil
}
