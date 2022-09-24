package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	c "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

// Config : Global variable to store the config
var Config *NwConfig

// NwConfig : stores all configuration of the application.
// The values are read from a config file or environment variable.
type NwConfig struct {

	// Database used to store NotionWatcher data's
	Database struct {
		DbType string `env:"DB_TYPE" validate:"required,oneof=sqlite3 mysql mariadb postgresql"`

		Sqlite struct {
			Dsn string `env:"DB_SQLITE_DSN" validate:"required_if=DbType sqlite3"`
		}

		Mysql struct {
			Host     string `env:"DB_MYSQL_HOST" validate:"required_if=DbType mysql"`
			Port     int    `env:"DB_MYSQL_PORT" validate:"required_if=DbType mysql"`
			Database string `env:"DB_MYSQL_DATABASE" validate:"required_if=DbType mysql"`
			User     string `env:"DB_MYSQL_USER" validate:"required_if=DbType mysql"`
			Password string `env:"DB_MYSQL_PASSWORD" validate:"required_if=DbType mysql"`
			Dsn      string `env:"DB_MYSQL_DSN"`
		}

		Mariadb struct {
			Host     string `env:"DB_MARIADB_HOST" validate:"required_if=DbType mariadb"`
			Port     int    `env:"DB_MARIADB_PORT" validate:"required_if=DbType mariadb"`
			Database string `env:"DB_MARIADB_DATABASE" validate:"required_if=DbType mariadb"`
			User     string `env:"DB_MARIADB_USER" validate:"required_if=DbType mariadb"`
			Password string `env:"DB_MARIADB_PASSWORD" validate:"required_if=DbType mariadb"`
			Dsn      string `env:"DB_MARIADB_DSN"`
		}

		Postgres struct {
			Host     string `env:"DB_POSTGRES_HOST" validate:"required_if=DbType postgresql"`
			Port     int    `env:"DB_POSTGRES_PORT" validate:"required_if=DbType postgresql"`
			Schema   string `env:"DB_POSTGRES_SCHEMA" validate:"required_if=DbType postgresql"`
			Database string `env:"DB_POSTGRES_DATABASE" validate:"required_if=DbType postgresql"`
			User     string `env:"DB_POSTGRES_USER" validate:"required_if=DbType postgresql"`
			Password string `env:"DB_POSTGRES_PASSWORD" validate:"required_if=DbType postgresql"`
			SslMode  string `env:"DB_POSTGRES_SSLMODE" validate:"required_if=DbType postgresql"`
			Dsn      string `env:"DB_POSTGRES_DSN"`
		}
	}

	// Environment
	// In development mode, logs are more verbose
	Environment string `env:"NW_ENV" validate:"required,oneof=development production"`

	// Logging
	LogPath  string `env:"LOG_PATH" validate:"required"`
	LogLevel string `env:"LOG_LEVEL" validate:"required,oneof=debug info warn error"`

	// Watchers
	// Dsn where json file are stored
	WatchersPath string `env:"WATCHER_PATH" validate:"required"`

	// Notion
	Token string `env:"NOTION_API_TOKEN" validate:"required"`
}

func init() {
	Config = &NwConfig{}

	// Set default config values
	Config.Database.DbType = "sqlite3"
	Config.Database.Sqlite.Dsn = "NotionWatcher.sqlite"
	Config.Database.Mysql.Host = "localhost"
	Config.Database.Mysql.Port = 3306
	Config.Database.Mariadb.Host = "localhost"
	Config.Database.Mariadb.Port = 3306
	Config.Database.Postgres.Host = "localhost"
	Config.Database.Postgres.Port = 5436
	Config.Database.Postgres.Schema = "public"
	Config.Database.Postgres.SslMode = "disable "
	Config.Environment = "production"
	Config.LogPath = "./logs/"
	Config.LogLevel = "info"
	Config.WatchersPath = "./watchers/"

	// Load configuration file
	Config.Load()
}

// Setup : Function called by golobby/config
func (nwConfig *NwConfig) Setup() {

	// Validate config data
	validate := validator.New()
	err := validate.Struct(nwConfig)
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Create connection string
	if nwConfig.Database.DbType == "mysql" {
		if strings.TrimSpace(nwConfig.Database.Mysql.Dsn) == "" {
			nwConfig.Database.Mysql.Dsn = fmt.Sprintf("Server=%s Port=%d Uid=%s Pwd=%s Database=%s",
				nwConfig.Database.Mysql.Host,
				nwConfig.Database.Mysql.Port,
				nwConfig.Database.Mysql.User,
				nwConfig.Database.Mysql.Password,
				nwConfig.Database.Mysql.Database)
		}
	} else if nwConfig.Database.DbType == "mariadb" {
		if strings.TrimSpace(nwConfig.Database.Mariadb.Dsn) == "" {
			nwConfig.Database.Mariadb.Dsn = fmt.Sprintf("Server=%s Port=%d Uid=%s Pwd=%s Database=%s",
				nwConfig.Database.Mariadb.Host,
				nwConfig.Database.Mariadb.Port,
				nwConfig.Database.Mariadb.User,
				nwConfig.Database.Mariadb.Password,
				nwConfig.Database.Mariadb.Database)
		}
	} else if nwConfig.Database.DbType == "postgresql" {
		if strings.TrimSpace(nwConfig.Database.Postgres.Dsn) == "" {
			nwConfig.Database.Postgres.Dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				nwConfig.Database.Postgres.Host,
				nwConfig.Database.Postgres.Port,
				nwConfig.Database.Postgres.User,
				nwConfig.Database.Postgres.Password,
				nwConfig.Database.Postgres.Database,
				nwConfig.Database.Postgres.SslMode)
		}
	}
}

// Load : Read configuration from file or environment variables.
func (nwConfig *NwConfig) Load() {
	// Define config feeder
	feeder1 := feeder.DotEnv{Path: ".env"}
	feeder2 := feeder.Env{}

	// Read config file
	err := c.New().AddFeeder(feeder1, feeder2).AddStruct(nwConfig).Feed()
	if err != nil {
		log.Fatal("cannot read config file:", err)
	}
}
