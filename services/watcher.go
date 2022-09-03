package services

import (
	"fmt"
	"time"

	"github.com/gennesseaux/NotionWatcher/models/json"
	nwLogger "github.com/gennesseaux/NotionWatcher/setup/logger"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

// logger : logger
var logger = nwLogger.Logger

type WatcherService struct {
	watchers []*json.Watcher
}

func NewWatcherService(watchers []*json.Watcher) *WatcherService {
	return &WatcherService{watchers: watchers}
}

func (ws *WatcherService) Run() {
	logger.Info("Running all watcher ...")
	for _, watcher := range ws.watchers {
		err := watcher.Run()
		if err != nil {
			logger.Info(fmt.Sprintf("Error running watcher : %s", watcher.Name), zap.Error(err))
		}
	}
}

func (ws *WatcherService) Watch() {
	logger.Info("Starting all watcher ...")

	s := gocron.NewScheduler(time.UTC)

	for _, watcher := range ws.watchers {
		err := watcher.Watch(s)
		if err != nil {
			logger.Info(fmt.Sprintf("Error running watcher : %s", watcher.Name), zap.Error(err))
		}
	}

	for ok := true; ok; ok = s.IsRunning() {
		time.Sleep(1 * time.Second)
	}
}
