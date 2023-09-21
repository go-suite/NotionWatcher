package watchers

import (
	"errors"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	log "github.com/go-mods/zerolog-quick"
	"github.com/go-playground/validator/v10"
	c "github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"os"
	"path/filepath"
	"regexp"
)

// config : config
var config = nwConfig.Config

type Watchers struct {
	Watchers []*watcher.Watcher
}

var Nw *Watchers

func init() {

	// Instance of NwWatcher
	Nw = &Watchers{}

	// Watchers folder
	if _, err := os.Stat(config.WatchersPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(config.WatchersPath, os.ModePerm)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot create watcher folder")
		}
	}

	// Loads existing watchers
	Nw.Load()
}

func (nw *Watchers) Load() {
	// List all json files in watcher folder
	var files []string
	err := filepath.Walk(config.WatchersPath, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".json", f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("loading watcher files")
	}

	// Convert json files to Watcher object
	for _, file := range files {
		// Create an instance of a watcher struct
		w := watcher.Watcher{}
		// Unmarshal
		err := c.New().AddFeeder(feeder.Json{Path: filepath.Join(config.WatchersPath, file)}).AddStruct(&w).Feed()
		if err != nil {
			log.Fatal().Err(err).Msgf("cannot read watcher file : '%s'", file)
		} else {
			// Validate watcher
			validate := validator.New()
			err := validate.Struct(w)
			if err != nil {
				log.Fatal().Err(err).Msgf("cannot read watcher file : '%s'", file)
			} else {
				// Add the watcher to the array
				nw.Watchers = append(nw.Watchers, &w)
			}
		}
	}
}

func (nw *Watchers) Get(name string) *watcher.Watcher {
	for _, w := range nw.Watchers {
		if w.Name == name {
			return w
		}
	}
	return nil
}

func (nw *Watchers) Add(watcher *watcher.Watcher) {
	nw.Watchers = append(nw.Watchers, watcher)
}

func (nw *Watchers) Remove(watcher *watcher.Watcher) {
	for i, w := range nw.Watchers {
		if w.Name == watcher.Name {
			nw.Watchers = append(nw.Watchers[:i], nw.Watchers[i+1:]...)
			return
		}
	}
}
