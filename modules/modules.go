package modules

import (
	// List of modules to import at startup
	_ "github.com/gennesseaux/NotionWatcher/modules/config"
	_ "github.com/gennesseaux/NotionWatcher/modules/datasource"
	_ "github.com/gennesseaux/NotionWatcher/modules/watchers"
)

func init() {
}
