package cmd

import (
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
)

type validateOptions struct {
	config *nwConfig.NwConfig
}

func newValidateCmd() *cobra.Command {

	// Options
	o := validateOptions{
		config: nwConfig.Config,
	}

	// Command
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the watcher",
		Args:  cobra.ExactArgs(0),
		RunE:  o.validateCmd,
	}

	// Disable usage print on error
	validateCmd.SilenceUsage = true

	return validateCmd
}

func (o *validateOptions) validateCmd(_ *cobra.Command, args []string) (err error) {

	// Get the watcher or return an error if it does not exist
	w := watchers.Get(args[0])
	if w == nil {
		return fmt.Errorf("the watcher %s does not exist", args[0])
	}

	// Validate the watcher and return an error if it's not valid
	if err = w.Validate(); err != nil {
		return fmt.Errorf("failed to validate the watcher: %s", err.Error())
	}

	log.Info().Msg("Watcher validated successfully!")

	return nil
}
