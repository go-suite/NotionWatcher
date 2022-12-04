package setup

import (
	// Call implicit init methods
	_ "github.com/gennesseaux/NotionWatcher/setup/config"
	_ "github.com/gennesseaux/NotionWatcher/setup/datasource"
	_ "github.com/gennesseaux/NotionWatcher/setup/watchers"
)

func init() {
}
