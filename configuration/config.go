package configuration

import (
	"log"
	"os"
)

const (
	portKey    = "AFTERHEE_PORT"
	dbPathKey  = "AFTERHEE_DUCKDB_FILENAME"
)

type Configuration struct {
	Port       string
	DBPath     string
}

func GetConfiguration() Configuration {
	return Configuration{
		Port:       getEnv(portKey, "8080"),
		DBPath:     getEnv(dbPathKey, "database/db.duckdb"),
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
