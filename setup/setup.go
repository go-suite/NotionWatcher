package setup

import (
	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"
	nwDatasource "github.com/gennesseaux/NotionWatcher/setup/datasource"
	nwLogger "github.com/gennesseaux/NotionWatcher/setup/logger"
	nwWatcher "github.com/gennesseaux/NotionWatcher/setup/watcher"

	// Call implicit init methods
	_ "github.com/gennesseaux/NotionWatcher/setup/config"
	_ "github.com/gennesseaux/NotionWatcher/setup/datasource"
	_ "github.com/gennesseaux/NotionWatcher/setup/logger"
	_ "github.com/gennesseaux/NotionWatcher/setup/watcher"
)

// logger : instance of Logger
var logger = nwLogger.Logger

// config : instance of Config
var config = nwConfig.Config

// datasource : instance of Datasource
var datasource = nwDatasource.Datasource

// watcher : instance of Watcher
var watcher = nwWatcher.Watcher

func init() {
	if config != nil {
		logger.Info("--> Config initialised")
	}
	if logger != nil {
		logger.Info("--> Logger initialised")
	}
	if datasource != nil {
		logger.Info("--> Datasource initialised")
	}
	if watcher != nil {
		logger.Info("--> Watcher initialised")
	}
}
