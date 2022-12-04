package cmd

import (
	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"
	"github.com/go-co-op/gocron"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
	"time"
)

type watchOptions struct {
	config *nwConfig.NwConfig
}

func newWatchCmd() *cobra.Command {

	// Options
	o := watchOptions{
		config: nwConfig.Config,
	}

	// Command
	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch all watchers",
		Run:   o.watchCmd,
	}

	// Disable usage print on error
	watchCmd.SilenceUsage = true

	return watchCmd
}

func (o *watchOptions) watchCmd(_ *cobra.Command, _ []string) {

	log.Info().Msg("Starting notion watcher ...")

	s := gocron.NewScheduler(time.UTC)

	for _, w := range watchers.Watchers {
		if w.Inactive {
			continue
		}
		if len(w.Cron) == 0 {
			continue
		}
		err := w.Watch(s)
		if err != nil {
			_ = s.RemoveByTag(w.Name)
		}
	}

	if o.config.Database.DbType == nwConfig.Sqlite3 {
		s.StartBlocking()
	} else {
		s.StartAsync()
	}

	for ok := true; ok; ok = s.IsRunning() {
		time.Sleep(1 * time.Second)
	}

	log.Info().Msg("Notion watcher exited!")
}
