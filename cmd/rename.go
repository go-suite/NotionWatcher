package cmd

import (
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
)

type renameOptions struct {
	config *nwConfig.NwConfig

	OldName string
	NewName string
}

func newRenameCmd() *cobra.Command {

	// Options
	o := renameOptions{
		config: nwConfig.Config,
	}

	// Command
	renameCmd := &cobra.Command{
		Use:     "rename [OldName] [NewName]",
		Aliases: []string{"r"},
		Short:   "Rename a watcher",
		Args:    cobra.ExactArgs(2),
		RunE:    o.renameCmd,
	}

	// Disable usage print on error
	renameCmd.SilenceUsage = true

	return renameCmd
}

func (o *renameOptions) renameCmd(_ *cobra.Command, args []string) (err error) {

	// Store the old name and the new name
	o.OldName = args[0]
	o.NewName = args[1]

	// Get the watcher or return an error
	ow := watchers.Get(o.OldName)
	if ow == nil {
		return fmt.Errorf("the watcher %s does not exist", args[0])
	}

	// Check if a watcher with the same name exist
	nw := watchers.Get(o.NewName)
	if nw != nil {
		return fmt.Errorf("a watcher with the name %s already exist", args[1])
	}

	// Create a copy of the watcher
	nw = new(watcher.Watcher)
	*nw = *ow
	nw.Name = o.NewName

	// Create the watcher and return an error if it fails
	err = nw.Create()
	if err != nil {
		return fmt.Errorf("failed to create the watcher: %s", err.Error())
	}

	// Delete the old watcher and return an error if it fails
	err = ow.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete the watcher: %s", err.Error())
	}

	// Remove the old watcher from the list of watchers
	watchers.Remove(ow)

	// Add the new watcher to the list
	watchers.Add(nw)

	log.Info().Msg("Watcher renamed successfully!")

	return nil
}
