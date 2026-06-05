package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"ping-uptime/internal/pkg/jwt"

	"github.com/joho/godotenv"
)

var (
	flagValues = make(map[string]string)
	flagOnce   sync.Once
)

// Initialize loads the .env file, parses command-line flags, and merges configuration.
func Initialize() error {
	if err := godotenv.Load(); err != nil {
		// .env file is optional — real env vars still work
		log.Println("No .env file found, reading from environment variables")
	}

	flagOnce.Do(func() {
		// Define flags with defaults from environment
		defaultAppName := os.Getenv("APP_NAME")
		if defaultAppName == "" {
			defaultAppName = "ping-uptime"
		}
		defaultServerMode := os.Getenv("SERVER_MODE")
		if defaultServerMode == "" {
			defaultServerMode = "info"
		}
		defaultPort := os.Getenv("PORT")
		if defaultPort == "" {
			defaultPort = "8080"
		}
		defaultHTTPTimeout := os.Getenv("HTTP_TIMEOUT")
		if defaultHTTPTimeout == "" {
			defaultHTTPTimeout = "60"
		}
		defaultCacheExpired := os.Getenv("CACHE_EXPIRED")
		if defaultCacheExpired == "" {
			defaultCacheExpired = "24"
		}
		defaultCachePurged := os.Getenv("CACHE_PURGED")
		if defaultCachePurged == "" {
			defaultCachePurged = "60"
		}
		defaultDBName := os.Getenv("DB_NAME")
		if defaultDBName == "" {
			defaultDBName = "ping-uptime"
		}
		defaultPoolConnLifetime := os.Getenv("POOL_CONN_LIFETIME")
		if defaultPoolConnLifetime == "" {
			defaultPoolConnLifetime = "60"
		}
		defaultJWTDayExpired := os.Getenv("JWT_DAY_EXPIRED")
		if defaultJWTDayExpired == "" {
			defaultJWTDayExpired = "60"
		}
		defaultJWTSignatureKey := os.Getenv("JWT_SIGNATURE_KEY")
		if defaultJWTSignatureKey == "" {
			defaultJWTSignatureKey = "4WSRLWxJdm"
		}

		_ = flag.String("app_name", defaultAppName, "Application Name")
		_ = flag.String("server_mode", defaultServerMode, "Server Mode (debug, info, etc.)")
		_ = flag.String("port", defaultPort, "Port to run the server on")
		_ = flag.String("http_timeout", defaultHTTPTimeout, "HTTP Timeout in seconds")
		_ = flag.String("cache_expired", defaultCacheExpired, "Cache expired in hours")
		_ = flag.String("cache_purged", defaultCachePurged, "Cache purged in minutes")
		_ = flag.String("db_name", defaultDBName, "Database Name")
		_ = flag.String("pool_conn_lifetime", defaultPoolConnLifetime, "Pool Connection Lifetime in minutes")
		_ = flag.String("jwt_day_expired", defaultJWTDayExpired, "JWT expiration time in days")
		_ = flag.String("jwt_signature_key", defaultJWTSignatureKey, "JWT signature key")

		if !flag.Parsed() {
			flag.Parse()
		}

		// Store only explicitly set flags
		flag.Visit(func(f *flag.Flag) {
			flagValues[f.Name] = f.Value.String()
		})
	})

	return nil
}

// checkKey aborts if the environment variable is not set.
func checkKey(key string) {
	if os.Getenv(key) == "" {
		log.Fatalf("Configuration key %s not found; aborting\n", key)
		os.Exit(1)
	}
}

// getStringNoPanic fetches configuration value from command-line flags first, then environment.
func getStringNoPanic(key string) string {
	flagKey := strings.ToLower(key)
	if val, ok := flagValues[flagKey]; ok {
		return val
	}
	return os.Getenv(key)
}

// GetString returns the value of configuration as a string, with a fallback default.
func GetString(key string, defaultValue ...string) string {
	val := getStringNoPanic(key)
	if val == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		checkKey(key)
	}
	return val
}

// GetInt returns the value of configuration as an int, with a fallback default.
func GetInt(key string, defaultValue ...int) int {
	valStr := getStringNoPanic(key)
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

// GetBool returns the value of configuration as a bool, with a fallback default.
func GetBool(key string, defaultValue ...bool) bool {
	valStr := getStringNoPanic(key)
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

// GetJWTService returns a JWT service using the signature key from the configuration.
func GetJWTService() jwt.JWT {
	signatureKey := GetString("JWT_SIGNATURE_KEY", "4WSRLWxJdm")
	jwtExpiredDays := GetInt("JWT_DAY_EXPIRED", 60)
	return jwt.NewJWTImpl(signatureKey, jwtExpiredDays)
}
