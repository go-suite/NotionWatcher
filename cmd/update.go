package cmd

import (
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
	"time"
)

type updateOptions struct {
	config *nwConfig.NwConfig

	Type        string
	DatabaseId  string
	StartDate   string
	DateLayout  string
	WebHook     string
	WebHookTest string
	Cron        string
	Token       string
	Inactive    bool
}

func newUpdateCmd() *cobra.Command {

	// Options
	o := updateOptions{
		config: nwConfig.Config,
	}

	// Command
	updateCmd := &cobra.Command{
		Use:     "update [Name]",
		Aliases: []string{"u"},
		Short:   "Update an existing watcher",
		Args:    cobra.ExactArgs(1),
		RunE:    o.updateCmd,
	}

	// Disable usage print on error
	updateCmd.SilenceUsage = true

	// Command flags
	updateCmd.Flags().StringVarP(&o.Type, "type", "", "", "type of event. Can be pageAddedToDatabase or pageUpdatedInDatabase")
	updateCmd.Flags().StringVarP(&o.DatabaseId, "database-id", "", "", "id of the database")
	updateCmd.Flags().StringVarP(&o.StartDate, "start-date", "", "", "minimum date to watch")
	updateCmd.Flags().StringVarP(&o.DateLayout, "date-layout", "", "2006-01-02", "layout to use when parsing the start date (default Y-m-d)")
	updateCmd.Flags().StringVarP(&o.WebHook, "hook", "", "", "WebHook to call on update")
	updateCmd.Flags().StringVarP(&o.WebHookTest, "hook-test", "", "", "WebHook to call for testing")
	updateCmd.Flags().StringVarP(&o.Cron, "cron", "", "", "cron")
	updateCmd.Flags().StringVarP(&o.Token, "token", "", "", "Notion token")
	updateCmd.Flags().BoolVarP(&o.Inactive, "inactive", "", false, "Inactive the watcher")

	return updateCmd
}

func (o *updateOptions) updateCmd(cmd *cobra.Command, args []string) (err error) {

	// Get the watcher or return an error
	w := watchers.Get(args[0])
	if w == nil {
		return fmt.Errorf("the watcher %s does not exist", args[0])
	}

	// Parse the type if provided
	if cmd.Flags().Changed("type") {
		eType, _ := event.ParseType(o.Type)
		w.Type = eType
	}

	// Parse the database id if provided
	if cmd.Flags().Changed("database-id") {
		w.DatabaseId = o.DatabaseId
	}

	// Parse the start date if provided
	if cmd.Flags().Changed("start-date") {
		startDate, err := time.Parse(o.DateLayout, o.StartDate)
		if err != nil {
			startDate = time.Now()
			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		}
		w.StartDate = &startDate
	}

	// Parse the web hook if provided
	if cmd.Flags().Changed("hook") {
		w.WebHook = o.WebHook
	}

	// Parse the web hook test if provided
	if cmd.Flags().Changed("hook-test") {
		w.WebHookTest = o.WebHookTest
	}

	// Parse the cron if provided
	if cmd.Flags().Changed("cron") {
		w.Cron = o.Cron
	}

	// Parse the token if provided
	if cmd.Flags().Changed("token") {
		w.Token = o.Token
	}

	// Parse the active state  if provided
	if cmd.Flags().Changed("inactive") {
		w.Inactive = o.Inactive
	}

	// Validate the watcher
	err = w.Validate()
	if err != nil {
		return err
	}

	// Update the watcher
	err = w.Create()
	if err != nil {
		return err
	}

	log.Info().Msg("Watcher updated successfully!")

	return nil
}
