package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHOST    string
	DBPORT    string
	DBUSER    string
	DBPASS    string
	DBNAME    string
	DBSSL     string
	APPPORT   string
	JWTSECRET string
}

func Load() *Config {

	// Load .env (safe ignore if not present)
	_ = godotenv.Load()

	cfg := &Config{
		DBHOST:    os.Getenv("DB_HOST"),
		DBPORT:    os.Getenv("DB_PORT"),
		DBUSER:    os.Getenv("DB_USER"),
		DBPASS:    os.Getenv("DB_PASS"),
		DBNAME:    os.Getenv("DB_NAME"),
		DBSSL:     os.Getenv("DB_SSL"),
		APPPORT:   os.Getenv("APP_PORT"),
		JWTSECRET: os.Getenv("JWT_SECRET"),
	}

	// -------------------------
	// Defaults
	// -------------------------
	if cfg.DBSSL == "" {
		cfg.DBSSL = "disable"
	}

	if cfg.APPPORT == "" {
		cfg.APPPORT = "8080"
	}

	if cfg.JWTSECRET == "" {
		log.Println("WARNING: JWT_SECRET not set → using default (dev only)")
		cfg.JWTSECRET = "CHANGE_ME_JWT_SECRET"
	}

	// -------------------------
	// Validation (important)
	// -------------------------
	if cfg.DBHOST == "" ||
		cfg.DBPORT == "" ||
		cfg.DBUSER == "" ||
		cfg.DBPASS == "" ||
		cfg.DBNAME == "" {

		log.Fatal("Missing required DB env: DB_HOST DB_PORT DB_USER DB_PASS DB_NAME")
	}

	return cfg
}
