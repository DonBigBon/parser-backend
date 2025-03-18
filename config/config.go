package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

type Config struct {
	DBServer   string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables")
	}

	config := &Config{
		DBServer:   getEnv("DB_SERVER", "localhost"),
		DBUser:     getEnv("DB_USER", "parser_user"),
		DBPassword: getEnv("DB_PASSWORD", "123123"),
		DBName:     getEnv("DB_NAME", "ParserDB"),
		DBPort:     getEnv("DB_PORT", "1433"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ConnectDB(config *Config) (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		config.DBServer, config.DBUser, config.DBPassword, config.DBPort, config.DBName)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
