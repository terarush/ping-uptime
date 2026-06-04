package config

import (
	"ping-uptime/internal/pkg/jwt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Initialize loads the .env file and environment variables.
func Initialize() error {
	if err := godotenv.Load(); err != nil {
		// .env file is optional — real env vars still work
		log.Println("No .env file found, reading from environment variables")
	}
	return nil
}

// checkKey aborts if the environment variable is not set.
func checkKey(key string) {
	if os.Getenv(key) == "" {
		log.Fatalf("Configuration key %s not found; aborting\n", key)
		os.Exit(1)
	}
}

// GetString returns the value of an environment variable as a string, with a fallback default.
func GetString(key string, defaultValue ...string) string {
	val := os.Getenv(key)
	if val == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		checkKey(key)
	}
	return val
}

// GetInt returns the value of an environment variable as an int, with a fallback default.
func GetInt(key string, defaultValue ...int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		checkKey(key)
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Configuration key %s is not a valid integer: %v\n", key, err)
		os.Exit(1)
	}
	return val
}

// GetBool returns the value of an environment variable as a bool, with a fallback default.
func GetBool(key string, defaultValue ...bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		checkKey(key)
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		log.Fatalf("Configuration key %s is not a valid boolean: %v\n", key, err)
		os.Exit(1)
	}
	return val
}

// GetJWTService returns a JWT service using the signature key from the environment.
func GetJWTService() jwt.JWT {
	signatureKey := GetString("JWT_SIGNATURE_KEY", "4WSRLWxJdm")
	jwtExpiredDays := GetInt("JWT_DAY_EXPIRED", 60)
	return jwt.NewJWTImpl(signatureKey, jwtExpiredDays)
}
