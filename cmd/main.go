package cmd

import (
	"fmt"
	"github.com/gennesseaux/NotionWatcher/modules/version"
	nwWatcher "github.com/gennesseaux/NotionWatcher/modules/watchers"
	"github.com/spf13/cobra"
)

// watchers : instance of Watcher
var watchers = nwWatcher.Nw

// Options of the main command
type mainOptions struct {
}

// main command
func newMainCmd() *cobra.Command {

	// Options
	o := mainOptions{}

	// Create the root command
	mainCmd := &cobra.Command{
		Use:     "notion-watcher",
		Aliases: []string{"nw"},
		Short:   "NotionWatcher give the ability to watch a Notion database",
		RunE:    o.mainCmd,
	}

	mainCmd.SetVersionTemplate(`customized version: {{.Version}}`)

	// version flag
	mainCmd.Flags().BoolP("version", "v", false, "Print the version number of NotionWatcher")

	// Sub commands
	mainCmd.AddCommand(newAddCmd())
	mainCmd.AddCommand(newUpdateCmd())
	mainCmd.AddCommand(newRenameCmd())
	mainCmd.AddCommand(newRemoveCmd())
	mainCmd.AddCommand(newValidateCmd())
	mainCmd.AddCommand(newListCmd())
	mainCmd.AddCommand(newInfoCmd())
	mainCmd.AddCommand(newRunCmd())
	mainCmd.AddCommand(newQuickRunCmd())
	mainCmd.AddCommand(newWatchCmd())

	//
	mainCmd.CompletionOptions.DisableDefaultCmd = true

	return mainCmd
}

//goland:noinspection GoBoolExpressions
func (o *mainOptions) mainCmd(cmd *cobra.Command, args []string) (err error) {

	// display version
	if v, _ := cmd.Flags().GetBool("version"); v {
		if version.Version != "" {
			if version.BuildDate == "" {
				cmd.Println(fmt.Sprintf("%s version %s", cmd.Name(), version.Version))
				return
			}
			if version.BuildDate == "" {
				cmd.Println(fmt.Sprintf("%s version %s, built on %s", cmd.Name(), version.Version, version.BuildDate))
				return
			}
		}
		return
	}

	// display help
	cmd.HelpFunc()(cmd, args)

	return
}

func Execute() (err error) {

	// Construct the root command
	mainCmd := newMainCmd()

	// Execute the root command
	if err = mainCmd.Execute(); err != nil {
		return
	}

	return nil
}
