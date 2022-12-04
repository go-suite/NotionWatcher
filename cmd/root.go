package cmd

import (
	"github.com/gennesseaux/NotionWatcher/setup/version"
	nwWatcher "github.com/gennesseaux/NotionWatcher/setup/watchers"
	"github.com/spf13/cobra"
)

// watchers : instance of Watcher
var watchers = nwWatcher.Nw

// Options of the root command
type rootOptions struct {
}

// root command
func newRootCmd() *cobra.Command {

	// Options
	o := rootOptions{}

	// Create the root command
	rootCmd := &cobra.Command{
		Use:     "nw",
		Short:   "NotionWatcher give the ability to watch a Notion database",
		Version: version.Version,
		RunE:    o.rootCmd,
	}

	// Sub commands
	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newUpdateCmd())
	rootCmd.AddCommand(newDeleteCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newWatchCmd())

	//
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}

func (o *rootOptions) rootCmd(cmd *cobra.Command, args []string) (err error) {

	// display help
	cmd.HelpFunc()(cmd, args)

	return
}

func Execute() (err error) {

	// Construct the root command
	rootCmd := newRootCmd()

	// Execute the root command
	if err = rootCmd.Execute(); err != nil {
		return
	}

	return nil
}
