package cmd

import (
	"errors"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/spf13/cobra"
)

type runOptions struct {
	config *nwConfig.NwConfig

	WebHook string
}

func newRunCmd() *cobra.Command {

	// Options
	o := runOptions{
		config: nwConfig.Config,
	}

	// Command
	runCmd := &cobra.Command{
		Use:   "run [Name]",
		Short: "Run a new watcher",
		Args:  cobra.ExactArgs(1),
		RunE:  o.runCmd,
	}
	//
	runCmd.SilenceUsage = true

	//
	runCmd.Flags().StringVarP(&o.WebHook, "hook", "", "", "WebHook to call on update")

	return runCmd
}

func (o *runOptions) runCmd(cmd *cobra.Command, args []string) (err error) {

	// Get the watcher
	w := watchers.Get(args[0])
	if w == nil {
		return errors.New("this watcher does not exist")
	}

	//
	if cmd.Flags().Changed("hook") {
		w.WebHook = o.WebHook
	}

	// Validate the watcher
	err = w.Validate()
	if w == nil {
		return err
	}

	// Run the watcher
	return w.Run()
}
