package env

import (
	"log"
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		log.Printf("\nReturning fallback for env var: %s\n", key)
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("\nReturning fallback for env var: %s\n", key)
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("\nReturning fallback for env var: %s\n", key)
		return fallback
	}

	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("\nReturning fallback for env var: %s\n", key)
		return fallback
	}

	valAsInt, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("\nReturning fallback for env var: %s\n", key)
		return fallback
	}

	return valAsInt
}
