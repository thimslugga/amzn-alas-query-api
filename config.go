package main

import (
	"log"
	"os"
	"strconv"
)

const (
	listenAddrEnvVar    = "LISTEN_ADDR"
	cacheTTLEnvVar      = "CACHE_TTL"
	redisHostEnvVar     = "REDIS_HOST"
	redisPasswordEnvVar = "REDIS_PASSWORD"
	redisDatabaseEnvVar = "REDIS_DATABASE"

	defaultListenAddr    = ":8080"
	defaultCacheTTL      = 300
	defaultRedisHost     = "redis:6379"
	defaultRedisPassword = ""
	defaultRedisDatabase = 0
)

// Global config
var config Config

// Config represents the application configuration
type Config struct {
	ListenAddr    string
	RedisHost     string
	RedisPassword string
	RedisDatabase int
	CacheTTL      int
}

// NewConfig parses the configuration from the environment and sets defaults
// for values not provided.
func NewConfig() (config Config) {
	config = Config{}

	// Check the cache TTL configuration
	cacheTTL := os.Getenv(cacheTTLEnvVar)
	if cacheTTL == "" {
		log.Println("Using default CacheTTL:", defaultCacheTTL)
		config.CacheTTL = defaultCacheTTL
	} else {
		parsed, err := strconv.Atoi(cacheTTL)
		if err != nil {
			log.Println("ERROR:", cacheTTL, "is not a valid integer, using default", defaultCacheTTL, "seconds")
			config.CacheTTL = defaultCacheTTL
		} else {
			log.Println("CacheTTL set to", parsed, "seconds")
			config.CacheTTL = parsed
		}
	}

	// Check the RedisDatabase configuration
	redisDB := os.Getenv(redisDatabaseEnvVar)
	if redisDB == "" {
		log.Println("Using default Redis Database:", defaultRedisDatabase)
		config.RedisDatabase = defaultRedisDatabase
	} else {
		parsed, err := strconv.Atoi(redisDB)
		if err != nil {
			log.Println("ERROR:", redisDB, "is not a valid integer, using default db", defaultRedisDatabase)
			config.RedisDatabase = defaultRedisDatabase
		} else {
			log.Println("RedisDatabase set to", parsed)
			config.RedisDatabase = parsed
		}
	}

	// Check the RedisPassword configuration
	redisPassword := os.Getenv(redisPasswordEnvVar)
	if redisPassword == "" {
		log.Println("Using default noauth for Redis connection")
		config.RedisPassword = defaultRedisPassword
	} else {
		log.Println("Parsed Redis Password from Environment")
		config.RedisPassword = redisPassword
	}

	// Check the RedisHost configuration
	redisHost := os.Getenv(redisHostEnvVar)
	if redisHost == "" {
		log.Println("Using default Redis Host:", defaultRedisHost)
		config.RedisHost = defaultRedisHost
	} else {
		log.Println("Connecting to Redis at", redisHost)
		config.RedisHost = redisHost
	}

	// Check the listener configuration
	listenAddr := os.Getenv(listenAddrEnvVar)
	if listenAddr == "" {
		log.Println("Using default listen address:", defaultListenAddr)
		config.ListenAddr = defaultListenAddr
	} else {
		log.Println("Listening for requests on", listenAddr)
		config.ListenAddr = listenAddr
	}

	return
}
