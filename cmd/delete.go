package cmd

import (
	"errors"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	config *nwConfig.NwConfig
}

func newDeleteCmd() *cobra.Command {

	// Options
	o := deleteOptions{
		config: nwConfig.Config,
	}

	// Command
	deleteCmd := &cobra.Command{
		Use:   "delete [Name]",
		Short: "Delete a watcher",
		Args:  cobra.ExactArgs(1),
		RunE:  o.deleteCmd,
	}

	return deleteCmd
}

func (o *deleteOptions) deleteCmd(cmd *cobra.Command, args []string) (err error) {

	// Get the watcher
	if watchers.Get(args[0]) == nil {
		return errors.New("this watcher does not exist")
	}

	w := &watcher.Watcher{
		Name: args[0],
	}

	// Delete the watcher
	err = w.Delete()
	if err != nil {
		return err
	}

	// Remove the watcher to the list of watchers
	watchers.Remove(w)

	log.Info().Msg("Watcher deleted successfully!")

	return nil
}
