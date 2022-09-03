package watcher

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gennesseaux/NotionWatcher/models/json"
	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"
	nwLogger "github.com/gennesseaux/NotionWatcher/setup/logger"

	"github.com/go-playground/validator/v10"
	c "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"go.uber.org/zap"
)

var Watcher *NwWatcher

// logger : logger
var logger = nwLogger.Logger

// config : config
var config = nwConfig.Config

type NwWatcher struct {
	Watchers []*json.Watcher
}

func init() {
	logger.Info("Loading watcher")

	// Instance of NwWatcher
	Watcher = &NwWatcher{}

	// Watchers folder
	if _, err := os.Stat(config.WatchersPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(config.WatchersPath, os.ModePerm)
		if err != nil {
			logger.Fatal("cannot create watcher folder:", zap.Error(err))
		}
	}

	// Load json files from watcher folder
	ext := ".json"
	var files []string
	err := filepath.Walk(config.WatchersPath, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	if err != nil {
		logger.Fatal("Error loading watcher file:", zap.Error(err))
	}

	// Convert json files to watcher object
	for _, file := range files {
		// Create an instance of a watcher struct
		watcher := json.Watcher{}
		// Unmarshal
		err := c.New().AddFeeder(feeder.Json{Path: filepath.Join(config.WatchersPath, file)}).AddStruct(&watcher).Feed()
		if err != nil {
			logger.Fatal("cannot read watcher file:", zap.Error(err))
		} else {
			// Validate watcher
			validate := validator.New()
			err := validate.Struct(watcher)
			if err != nil {
				logger.Fatal("cannot read watcher file:", zap.Error(err))
			} else {
				// Add the watcher to the array
				Watcher.Watchers = append(Watcher.Watchers, &watcher)
			}
		}
	}
}
