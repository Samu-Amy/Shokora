package env

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func GetString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		log.Printf("Returning fallback for env var: %s\n\n", key)
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Returning fallback for env var: %s\n\n", key)
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Returning fallback for env var: %s\n\n", key)
		return fallback
	}

	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Returning fallback for env var: %s\n\n", key)
		return fallback
	}

	valAsInt, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("Returning fallback for env var: %s\n\n", key)
		return fallback
	}

	return valAsInt
}

func LoadCORSOrigins(fallback []string) []string {
	originsEnv, ok := os.LookupEnv("ALLOWED_ORIGINS")
	if !ok || originsEnv == "" {
		log.Printf("Returning fallback for env var: ALLOWED_ORIGINS")
		return fallback
	}

	origins := strings.Split(originsEnv, ",")

	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return origins
}
