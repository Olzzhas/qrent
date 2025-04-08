package main

import (
	"os"
	"strconv"
)

func loadConfig() config {
	var cfg config

	cfg.port = getEnvAsInt("PORT", 4000)
	cfg.env = getEnv("ENV", "production")

	cfg.db.dsn = getEnv("DB_DSN", "postgres://username:password@localhost:5432/dbname?sslmode=disable")
	cfg.db.maxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	cfg.db.maxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", 100)
	cfg.db.maxIdleTime = getEnv("DB_MAX_IDLE_TIME", "15m")

	cfg.limiter.rps = getEnvAsFloat("LIMITER_RPS", 10)
	cfg.limiter.burst = getEnvAsInt("LIMITER_BURST", 50)
	cfg.limiter.enabled = getEnvAsBool("LIMITER_ENABLED", true)

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := getEnv(name, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsFloat(name string, defaultVal float64) float64 {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
