package datasource

import (
	"context"
	"database/sql"
	wdb "github.com/gennesseaux/NotionWatcher/common/db"
	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"
	loggorm "github.com/go-mods/zerolog-gorm"
	log "github.com/go-mods/zerolog-quick"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
		db, err = gorm.Open(sqlite.Open(config.Database.Sqlite.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Sqlite database")
		}
		sqlDb, _ = db.DB()
	} else if config.Database.DbType == nwConfig.Mysql {
		db, err = gorm.Open(mysql.Open(config.Database.Mysql.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Mysql database")
		}
		sqlDb, _ = db.DB()
		sqlDb.SetConnMaxLifetime(time.Minute * 3)
		sqlDb.SetMaxOpenConns(10)
		sqlDb.SetMaxIdleConns(10)
	} else if config.Database.DbType == nwConfig.Mariadb {
		db, err = gorm.Open(mysql.Open(config.Database.Mariadb.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Mariadb database")
		}
		sqlDb, _ = db.DB()
		sqlDb.SetConnMaxLifetime(time.Minute * 3)
		sqlDb.SetMaxOpenConns(10)
		sqlDb.SetMaxIdleConns(10)
	} else if config.Database.DbType == nwConfig.Postgresql {
		db, err = gorm.Open(postgres.Open(config.Database.Postgres.Dsn), &gormConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot connect to Postgresql database")
		}
	}

	//
	db = db.WithContext(gormLogger.WithContext(context.Background()))

	// Enable debug mode while in development
	if config.Environment == "development" {
		db = db.Debug()
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
		&wdb.Database{},
		&wdb.DatabaseWatcher{},
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
