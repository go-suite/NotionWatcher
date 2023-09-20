package cmd

import (
	"errors"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
	"time"
)

type updateOptions struct {
	config *nwConfig.NwConfig

	Type       string
	DatabaseId string
	StartDate  string
	DateLayout string
	WebHook    string
	Cron       string
	Token      string
	Inactive   bool
}

func newUpdateCmd() *cobra.Command {

	// Options
	o := updateOptions{
		config: nwConfig.Config,
	}

	// Command
	updateCmd := &cobra.Command{
		Use:   "update [Name]",
		Short: "Update a new watcher",
		Args:  cobra.ExactArgs(1),
		RunE:  o.updateCmd,
	}

	// Command flags
	updateCmd.Flags().StringVarP(&o.Type, "type", "", "", "type of event. Can be pageAddedToDatabase or pageUpdatedInDatabase")
	updateCmd.Flags().StringVarP(&o.DatabaseId, "database-id", "", "", "id of the database")
	updateCmd.Flags().StringVarP(&o.StartDate, "start-date", "", "", "minimum date to watch")
	updateCmd.Flags().StringVarP(&o.DateLayout, "date-layout", "", "2006-01-02", "layout to use when parsing the start date (default Y-m-d)")
	updateCmd.Flags().StringVarP(&o.WebHook, "hook", "", "", "WebHook to call on update")
	updateCmd.Flags().StringVarP(&o.Cron, "cron", "", "", "cron")
	updateCmd.Flags().StringVarP(&o.Token, "token", "", "", "Notion token")
	updateCmd.Flags().BoolVarP(&o.Inactive, "inactive", "", false, "Inactive the watcher")

	// Command required flags
	_ = updateCmd.MarkFlagRequired("type")
	_ = updateCmd.MarkFlagRequired("database-id")
	_ = updateCmd.MarkFlagRequired("hook")
	_ = updateCmd.MarkFlagRequired("token")

	return updateCmd
}

func (o *updateOptions) updateCmd(cmd *cobra.Command, args []string) (err error) {

	// Get the watcher
	w := watchers.Get(args[0])
	if w == nil {
		return errors.New("this watcher does not exist")
	}

	if cmd.Flags().Changed("type") {
		eType, _ := event.ParseType(o.Type)
		w.Type = eType
	}

	if cmd.Flags().Changed("database-id") {
		w.DatabaseId = o.DatabaseId
	}

	if cmd.Flags().Changed("start-date") {
		startDate, err := time.Parse(o.DateLayout, o.StartDate)
		if err != nil {
			startDate = time.Now()
			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		}
	}

	if cmd.Flags().Changed("hook") {
		w.WebHook = o.WebHook
	}

	if cmd.Flags().Changed("cron") {
		w.Cron = o.Cron
	}

	if cmd.Flags().Changed("token") {
		w.Token = o.Token
	}

	if cmd.Flags().Changed("inactive") {
		w.Inactive = o.Inactive
	}

	startDate, err := time.Parse(o.DateLayout, o.StartDate)
	if err != nil {
		startDate = time.Now()
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
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
