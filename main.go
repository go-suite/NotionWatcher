package main

import (
	"github.com/gennesseaux/NotionWatcher/services"
	nwWatcher "github.com/gennesseaux/NotionWatcher/setup/watcher"

	_ "github.com/gennesseaux/NotionWatcher/setup"
)

func main() {
	watchers := nwWatcher.Watcher
	watcherService := services.NewWatcherService(watchers.Watchers)
	watcherService.Run()
	watcherService.Watch()
}
