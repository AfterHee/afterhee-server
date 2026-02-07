package configuration

import (
	"log"
	"os"
)

const (
	portKey          = "AFTERHEE_PORT"
	dbPathKey        = "AFTERHEE_DUCKDB_FILENAME"
	redisHostKey     = "AFTERHEE_REDIS_HOST"
	redisPortKey     = "AFTERHEE_REDIS_PORT"
	redisPasswordKey = "AFTERHEE_REDIS_PASSWORD"
	neisAPIKey       = "AFTERHEE_NEIS_API_KEY"
	geminiAPIKey     = "AFTERHEE_GEMINI_API_KEY"
	runningProfile   = "AFTERHEE_PROFILE"
)

type Configuration struct {
	Port          string
	DBPath        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	NEISAPIKey    string
	GeminiAPIKey  string
	IsDevMode     bool
}

func GetConfiguration() Configuration {
	return Configuration{
		Port:          getEnv(portKey, "8080"),
		DBPath:        getEnv(dbPathKey, "database/db.duckdb"),
		RedisHost:     getEnv(redisHostKey, "localhost"),
		RedisPort:     getEnv(redisPortKey, "56379"),
		RedisPassword: getEnv(redisPasswordKey, ""),
		NEISAPIKey:    getEnv(neisAPIKey, ""),
		GeminiAPIKey:  getEnv(geminiAPIKey, ""),
		IsDevMode:     getEnv(runningProfile, "PRODUCTION") == "DEV",
	}
}

func getEnv(envKey string, fallback string) string {
	envValue := os.Getenv(envKey)

	if envValue == "" {
		if fallback == "" {
			log.Fatalln("environment value(key=" + envKey + ") is empty but fallback value is also empty.")
		}

		return fallback
	}

	return envValue
}
