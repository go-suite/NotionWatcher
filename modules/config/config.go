package config

import (
	"fmt"
	log "github.com/go-mods/zerolog-quick"
	"github.com/go-mods/zerolog-quick/console/colored"
	"github.com/go-playground/validator/v10"
	c "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

// Config : Global variable to store the config
var Config *NwConfig

// Database types
const (
	Sqlite3    string = "sqlite3"
	Mysql      string = "Mmsql"
	Mariadb    string = "mariadb"
	Postgresql string = "postgresql"
)

// NwConfig : stores the configuration of the application.
// The values are read from a config file or environment variable.
type NwConfig struct {

	// Database used to store NotionWatcher data's
	Database struct {
		DbType string `env:"DB_TYPE" validate:"required,oneof=sqlite3 mysql mariadb postgresql"`

		Sqlite struct {
			Dsn string `env:"DB_SQLITE_DSN" validate:"required_if=DbType Sqlite3"`
		}

		Mysql struct {
			Host     string `env:"DB_MYSQL_HOST" validate:"required_if=DbType Mysql"`
			Port     int    `env:"DB_MYSQL_PORT" validate:"required_if=DbType Mysql"`
			Database string `env:"DB_MYSQL_DATABASE" validate:"required_if=DbType Mysql"`
			User     string `env:"DB_MYSQL_USER" validate:"required_if=DbType Mysql"`
			Password string `env:"DB_MYSQL_PASSWORD" validate:"required_if=DbType Mysql"`
			Dsn      string `env:"DB_MYSQL_DSN"`
		}

		Mariadb struct {
			Host     string `env:"DB_MARIADB_HOST" validate:"required_if=DbType Mariadb"`
			Port     int    `env:"DB_MARIADB_PORT" validate:"required_if=DbType Mariadb"`
			Database string `env:"DB_MARIADB_DATABASE" validate:"required_if=DbType Mariadb"`
			User     string `env:"DB_MARIADB_USER" validate:"required_if=DbType Mariadb"`
			Password string `env:"DB_MARIADB_PASSWORD" validate:"required_if=DbType Mariadb"`
			Dsn      string `env:"DB_MARIADB_DSN"`
		}

		Postgres struct {
			Host     string `env:"DB_POSTGRES_HOST" validate:"required_if=DbType Postgresql"`
			Port     int    `env:"DB_POSTGRES_PORT" validate:"required_if=DbType Postgresql"`
			Schema   string `env:"DB_POSTGRES_SCHEMA" validate:"required_if=DbType Postgresql"`
			Database string `env:"DB_POSTGRES_DATABASE" validate:"required_if=DbType Postgresql"`
			User     string `env:"DB_POSTGRES_USER" validate:"required_if=DbType Postgresql"`
			Password string `env:"DB_POSTGRES_PASSWORD" validate:"required_if=DbType Postgresql"`
			SslMode  string `env:"DB_POSTGRES_SSLMODE" validate:"required_if=DbType Postgresql"`
			Dsn      string `env:"DB_POSTGRES_DSN"`
		}
	}

	// Environment
	// In development mode, logs are more verbose
	Environment string `env:"NW_ENV" validate:"required,oneof=development production"`

	// Logging
	// LogFile holds the configuration used for the log file
	LogFile  logFile
	LogLevel string `env:"LOG_LEVEL" validate:"required,oneof=debug info warn error"`

	// Watchers
	// Path where watcher files are stored
	WatchersPath string `env:"WATCHER_PATH" validate:"required"`
}

type logFile struct {
	// Filename is the file to write logs to. Backup log files will be retained
	// in the same directory. It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `env:"LOG_FILE" validate:"required"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `env:"LOG_MAX_SIZE"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `env:"LOG_MAX_AGE"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `env:"LOG_MAX_BACKUPS"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `env:"LOG_LOCAL_TIME"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `env:"LOG_COMPRESS"`
	// contains filtered or unexported fields
}

func init() {
	Config = &NwConfig{}

	// Set default config values
	Config.Database.DbType = Sqlite3
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
	Config.LogFile = logFile{
		Filename: "./logs/notion_watcher.log",
		MaxSize:  10,
	}
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
		log.Fatal().Err(err).Msg("loading .env file")
	}

	// Create connection string
	if nwConfig.Database.DbType == Mysql {
		if strings.TrimSpace(nwConfig.Database.Mysql.Dsn) == "" {
			nwConfig.Database.Mysql.Dsn = fmt.Sprintf("Server=%s Port=%d Uid=%s Pwd=%s Database=%s",
				nwConfig.Database.Mysql.Host,
				nwConfig.Database.Mysql.Port,
				nwConfig.Database.Mysql.User,
				nwConfig.Database.Mysql.Password,
				nwConfig.Database.Mysql.Database)
		}
	} else if nwConfig.Database.DbType == Mariadb {
		if strings.TrimSpace(nwConfig.Database.Mariadb.Dsn) == "" {
			nwConfig.Database.Mariadb.Dsn = fmt.Sprintf("Server=%s Port=%d Uid=%s Pwd=%s Database=%s",
				nwConfig.Database.Mariadb.Host,
				nwConfig.Database.Mariadb.Port,
				nwConfig.Database.Mariadb.User,
				nwConfig.Database.Mariadb.Password,
				nwConfig.Database.Mariadb.Database)
		}
	} else if nwConfig.Database.DbType == Postgresql {
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

	// default logger
	var logLevel = zerolog.InfoLevel
	if nwConfig.LogLevel == "debug" {
		logLevel = zerolog.DebugLevel
	} else if nwConfig.LogLevel == "info" {
		logLevel = zerolog.InfoLevel
	} else if nwConfig.LogLevel == "warn" {
		logLevel = zerolog.WarnLevel
	} else if nwConfig.LogLevel == "error" {
		logLevel = zerolog.ErrorLevel
	}
	log.Logger = colored.Message.Level(logLevel)
}

// Load : Read configuration from file or environment variables.
func (nwConfig *NwConfig) Load() {

	config := c.New()

	curDir, _ := os.Getwd()

	if _, err := os.Stat(".env"); err == nil {
		config.AddFeeder(feeder.DotEnv{Path: ".env"})
	} else {
		log.Debug().Msgf(".env file not found in %s", curDir)
	}

	config.AddFeeder(feeder.Env{})

	// Read config file
	err := config.AddStruct(nwConfig).Feed()
	if err != nil {
		log.Fatal().Err(err).Msg("reading .env file")
	} else {
		log.Debug().Msg("config loaded")
	}
}
