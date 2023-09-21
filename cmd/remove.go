package cmd

import (
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
)

type removeOptions struct {
	config *nwConfig.NwConfig
}

func newRemoveCmd() *cobra.Command {

	// Options
	o := removeOptions{
		config: nwConfig.Config,
	}

	// Command
	removeCmd := &cobra.Command{
		Use:     "remove [Name]",
		Aliases: []string{"r", "delete", "d"},
		Short:   "Remove a watcher",
		Args:    cobra.ExactArgs(1),
		RunE:    o.removeCmd,
	}

	// Disable usage print on error
	removeCmd.SilenceUsage = true

	return removeCmd
}

func (o *removeOptions) removeCmd(_ *cobra.Command, args []string) (err error) {

	// Get the watcher or return an error if it does not exist
	w := watchers.Get(args[0])
	if w == nil {
		return fmt.Errorf("the watcher %s does not exist", args[0])
	}

	// Delete the watcher and return an error if it fails
	err = w.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete the watcher %s: %s", args[0], err.Error())
	}

	// Remove the watcher from the list of watchers
	watchers.Remove(w)

	log.Info().Msg("Watcher removed successfully!")

	return nil
}
