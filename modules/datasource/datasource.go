package datasource

import (
	"context"
	"database/sql"
	"errors"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/database/models"
	loggorm "github.com/go-mods/zerolog-gorm"
	log "github.com/go-mods/zerolog-quick"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

// Datasource :
var Datasource *NwDatasource

// config : config
var config = nwConfig.Config

type NwDatasource struct {
	db    *gorm.DB
	sqlDb *sql.DB
}

func init() {

	var db *gorm.DB
	var sqlDb *sql.DB
	var err error

	// file writer
	fileWriter := &lumberjack.Logger{
		Filename:   config.LogFile.Filename,
		MaxSize:    config.LogFile.MaxSize,
		MaxBackups: config.LogFile.MaxBackups,
		MaxAge:     config.LogFile.MaxAge,
		Compress:   config.LogFile.Compress,
	}

	// gorm logger
	gormLogger := zerolog.New(fileWriter).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	// gorm config
	gormConfig := gorm.Config{Logger: &loggorm.GormLogger{}}

	if config.Database.DbType == nwConfig.Sqlite3 {
		// Get the database parent folder
		pf := filepath.Dir(config.Database.Sqlite.Dsn)
		// Create database file parent folder if not exists
		if _, err := os.Stat(pf); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(pf, os.ModePerm)
			if err != nil {
				log.Fatal().Err(err).Msg("cannot create database folder")
			}
		}
		// Open the database
		db, err = gorm.Open(sqlite.Open(config.Database.Sqlite.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Sqlite database")
		} else {
			log.Debug().Msgf("Opened sqlite database at %s", config.Database.Sqlite.Dsn)
		}
		sqlDb, _ = db.DB()
	} else if config.Database.DbType == nwConfig.Mysql {
		db, err = gorm.Open(mysql.Open(config.Database.Mysql.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Mysql database")
		} else {
			log.Debug().Msgf("Opened Mysql database at %s", config.Database.Mysql.Dsn)
		}
		sqlDb, _ = db.DB()
		sqlDb.SetConnMaxLifetime(time.Minute * 3)
		sqlDb.SetMaxOpenConns(10)
		sqlDb.SetMaxIdleConns(10)
	} else if config.Database.DbType == nwConfig.Mariadb {
		db, err = gorm.Open(mysql.Open(config.Database.Mariadb.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Mariadb database")
		} else {
			log.Debug().Msgf("Opened Mariadb database at %s", config.Database.Mariadb.Dsn)
		}
		sqlDb, _ = db.DB()
		sqlDb.SetConnMaxLifetime(time.Minute * 3)
		sqlDb.SetMaxOpenConns(10)
		sqlDb.SetMaxIdleConns(10)
	} else if config.Database.DbType == nwConfig.Postgresql {
		db, err = gorm.Open(postgres.Open(config.Database.Postgres.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Postgresql database")
		} else {
			log.Debug().Msgf("Opened Postgres database at %s", config.Database.Postgres.Dsn)
		}
	}

	//
	db = db.WithContext(gormLogger.WithContext(context.Background()))

	// Enable debug mode while in development
	if config.Environment == "development" {
		db = db.Debug()
		log.Debug().Msg("Environment set to development")
	}

	// try to establish connection
	if sqlDb != nil {
		err = sqlDb.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to database")
		}
	}

	// Migrate models to database
	err = db.AutoMigrate(
		&models.Database{},
		&models.DatabaseWatcher{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot migrate database")
	}

	Datasource = &NwDatasource{db: db, sqlDb: sqlDb}
}

func (d NwDatasource) DB() *gorm.DB {
	return d.db
}

func (d NwDatasource) SqlDB() *sql.DB {
	return d.sqlDb
}
