package datasource

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"moul.io/zapgorm2"
	"time"

	dbModels "github.com/gennesseaux/NotionWatcher/models/db"
	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"
	nwLogger "github.com/gennesseaux/NotionWatcher/setup/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Datasource :
var Datasource *NwDatasource

// logger : logger
var logger = nwLogger.Logger

// config : config
var config = nwConfig.Config

type NwDatasource struct {
	db    *gorm.DB
	sqlDb *sql.DB
}

func init() {
	logger.Info("Connecting to database")

	var db *gorm.DB
	var sqlDb *sql.DB
	var err error

	gormLogger := zapgorm2.New(logger)

	if config.Database.DbType == "sqlite3" {
		db, err = gorm.Open(sqlite.Open(config.Database.Sqlite.Dsn), &gorm.Config{Logger: gormLogger})
		if err != nil {
			logger.Fatal("cannot connect to db:", zap.Error(err))
		}
		sqlDb, _ = db.DB()
	} else if config.Database.DbType == "mysql" {
		db, err = gorm.Open(mysql.Open(config.Database.Mysql.Dsn), &gorm.Config{Logger: gormLogger})
		if err != nil {
			logger.Fatal("cannot connect to db:", zap.Error(err))
		}
		sqlDb, _ = db.DB()
		sqlDb.SetConnMaxLifetime(time.Minute * 3)
		sqlDb.SetMaxOpenConns(10)
		sqlDb.SetMaxIdleConns(10)
	} else if config.Database.DbType == "mariadb" {
		db, err = gorm.Open(mysql.Open(config.Database.Mariadb.Dsn), &gorm.Config{Logger: gormLogger})
		if err != nil {
			logger.Fatal("cannot connect to db:", zap.Error(err))
		}
		sqlDb, _ = db.DB()
		sqlDb.SetConnMaxLifetime(time.Minute * 3)
		sqlDb.SetMaxOpenConns(10)
		sqlDb.SetMaxIdleConns(10)
	} else if config.Database.DbType == "postgresql" {
		db, err = gorm.Open(postgres.Open(config.Database.Postgres.Dsn), &gorm.Config{Logger: gormLogger})
		if err != nil {
			logger.Fatal("cannot connect to db:", zap.Error(err))
		}
	}

	// Enable debug mode while in development
	if config.Environment == "development" {
		db = db.Debug()
	}

	// try to establish connection
	if sqlDb != nil {
		err = sqlDb.Ping()
		if err != nil {
			logger.Fatal("cannot connect to db:", zap.Error(err))
		}
	}

	// Migrate models to database
	err = db.AutoMigrate(
		&dbModels.Database{},
		&dbModels.DatabaseWatcher{},
	)
	if err != nil {
		logger.Error("cannot migrate database", zap.Error(err))
	}

	Datasource = &NwDatasource{db: db, sqlDb: sqlDb}
}

func (d NwDatasource) DB() *gorm.DB {
	return d.db
}

func (d NwDatasource) SqlDB() *sql.DB {
	return d.sqlDb
}
