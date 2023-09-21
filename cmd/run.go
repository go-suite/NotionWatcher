package cmd

import (
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/spf13/cobra"
)

type runOptions struct {
	config *nwConfig.NwConfig

	WebHook string
	Test    bool
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

	// Command flags
	runCmd.Flags().StringVarP(&o.WebHook, "hook", "", "", "WebHook to call on update")
	runCmd.Flags().BoolVarP(&o.Test, "test", "", false, "Use the test the WebHook")

	return runCmd
}

func (o *runOptions) runCmd(cmd *cobra.Command, args []string) (err error) {

	// Get the watcher or return an error
	w := watchers.Get(args[0])
	if w == nil {
		return fmt.Errorf("the watcher %s does not exist", args[0])
	}

	// Parse the WebHook if provided
	if cmd.Flags().Changed("hook") {
		w.WebHook = o.WebHook
	}

	// Validate the watcher and return an error if it's not valid
	if err = w.Validate(); err != nil {
		return fmt.Errorf("failed to validate the watcher: %s", err.Error())
	}

	// Run the watcher
	if o.Test {
		return w.RunTest()
	}
	return w.Run()
}
