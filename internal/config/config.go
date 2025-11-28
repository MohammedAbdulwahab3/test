package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBDriver    string // sqlite or postgres
	DBDSN       string
	AutoMigrate bool
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func Load() Config {
	_ = godotenv.Load()

	port := getenv("PORT", "8080")
	driver := getenv("DB_DRIVER", "sqlite")

	var dsn string
	if driver == "postgres" {
		dsn = getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/app?sslmode=disable")
	} else {
		// fast defaults for SQLite with WAL and busy timeout
		path := getenv("SQLITE_PATH", "data.db")
		dsn = getenv("DB_DSN", "file:"+path+"?cache=shared&_fk=1&_journal_mode=WAL&_busy_timeout=10000")
	}

	autoMigrate, _ := strconv.ParseBool(getenv("AUTO_MIGRATE", "true"))

	return Config{
		Port:        port,
		DBDriver:    driver,
		DBDSN:       dsn,
		AutoMigrate: autoMigrate,
	}
}
