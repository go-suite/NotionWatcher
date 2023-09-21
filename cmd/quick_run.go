package cmd

import (
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	"github.com/spf13/cobra"
)

type quickRunOptions struct {
	config *nwConfig.NwConfig

	Type       string
	DatabaseId string
	WebHook    string
	Token      string
}

func newQuickRunCmd() *cobra.Command {

	// Options
	o := quickRunOptions{
		config: nwConfig.Config,
	}

	// Command
	quickRunCmd := &cobra.Command{
		Use:     "quick-run [Name]",
		Aliases: []string{"qr"},
		Short:   "Quickly run a watcher",
		RunE:    o.runCmd,
	}

	//
	quickRunCmd.SilenceUsage = true

	// Command flags
	quickRunCmd.Flags().StringVarP(&o.Type, "type", "", "", "type of event. Can be pageAddedToDatabase or pageUpdatedInDatabase")
	quickRunCmd.Flags().StringVarP(&o.DatabaseId, "database-id", "", "", "id of the database")
	quickRunCmd.Flags().StringVarP(&o.WebHook, "hook", "", "", "WebHook to call on update")
	quickRunCmd.Flags().StringVarP(&o.Token, "token", "", "", "Notion token")

	return quickRunCmd
}

func (o *quickRunOptions) runCmd(_ *cobra.Command, _ []string) (err error) {

	// Create a new watcher based on parameters
	w := watcher.Watcher{
		Name:       "QuickWatcher",
		Type:       event.MustParseType(o.Type),
		DatabaseId: o.DatabaseId,
		WebHook:    o.WebHook,
		Token:      o.Token,
	}

	// Validate the watcher and return an error if it's not valid
	if err = w.Validate(); err != nil {
		return
	}

	// Run the watcher
	return w.Run()
}
