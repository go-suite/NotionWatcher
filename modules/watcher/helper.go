package watcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	"github.com/go-co-op/gocron"
	log "github.com/go-mods/zerolog-quick"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
)

func (w *Watcher) Create() error {
	file, _ := json.MarshalIndent(w, "", "  ")
	fileName := fmt.Sprintf("%s.json", w.Name)
	return os.WriteFile(filepath.Join(config.WatchersPath, fileName), file, 0644) // #nosec:G306
}

func (w *Watcher) Delete() error {
	fileName := fmt.Sprintf("%s.json", w.Name)
	return os.Remove(filepath.Join(config.WatchersPath, fileName))
}

func (w *Watcher) Run() (err error) {
	return w.run(false)
}

func (w *Watcher) RunTest() (err error) {
	return w.run(true)
}

func (w *Watcher) run(useTestHook bool) (err error) {
	// if in test mode, use the test hook
	if useTestHook {
		w.WebHook = w.WebHookTest
		log.Info().Msgf("running watcher: %s in test mode", w.Name)
	} else {
		log.Info().Msgf("running watcher: %s", w.Name)
	}

	if w.Type == event.PageAddedToDatabase || w.Type == event.PageUpdatedInDatabase {
		err = w.prepareDatabaseWatcher()
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to prepare database watcher: %s", w.Name))
			return err
		}
		err = w.runDatabaseWatcher()
		if err != nil {
			log.WithLevel(zerolog.ErrorLevel).Err(err).Send()
			return err
		}
		return nil
	}

	return errors.New("function not defined")
}

func (w *Watcher) Watch(s *gocron.Scheduler) (err error) {
	log.Info().Msgf("watcher: %s", w.Name)
	_, err = s.Cron(w.Cron).Tag(w.Name).Do(w.Run)
	if err != nil {
		return err
	}
	return
}
