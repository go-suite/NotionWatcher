package cmd

import (
	"errors"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
	"time"
)

type createOptions struct {
	config *nwConfig.NwConfig

	Name       string
	Type       string
	DatabaseId string
	StartDate  string
	DateLayout string
	WebHook    string
	Cron       string
	Token      string
	Inactive   bool
}

func newCreateCmd() *cobra.Command {

	// Options
	o := createOptions{
		config: nwConfig.Config,
	}

	// Command
	createCmd := &cobra.Command{
		Use:   "create [Name]",
		Short: "Create a new watcher",
		Args:  cobra.ExactArgs(1),
		RunE:  o.createCmd,
	}

	// Command flags
	createCmd.Flags().StringVarP(&o.Type, "type", "", "", "type of event. Can be pageAddedToDatabase or pageUpdatedInDatabase")
	createCmd.Flags().StringVarP(&o.DatabaseId, "database-id", "", "", "id of the database")
	createCmd.Flags().StringVarP(&o.StartDate, "start-date", "", "", "minimum date to watch")
	createCmd.Flags().StringVarP(&o.DateLayout, "date-layout", "", "2006-01-02", "layout to use when parsing the start date (default Y-m-d)")
	createCmd.Flags().StringVarP(&o.WebHook, "hook", "", "", "WebHook to call on update")
	createCmd.Flags().StringVarP(&o.Cron, "cron", "", "", "cron")
	createCmd.Flags().StringVarP(&o.Token, "token", "", "", "Notion token")
	createCmd.Flags().BoolVarP(&o.Inactive, "inactive", "", false, "Inactive the watcher")

	// Command required flags
	_ = createCmd.MarkFlagRequired("type")
	_ = createCmd.MarkFlagRequired("database-id")
	_ = createCmd.MarkFlagRequired("hook")
	_ = createCmd.MarkFlagRequired("token")

	return createCmd
}

func (o *createOptions) createCmd(cmd *cobra.Command, args []string) (err error) {

	// Check if a watcher with the same name exist
	if watchers.Get(args[0]) != nil {
		return errors.New("this watcher already exist")
	}

	eType, _ := event.ParseType(o.Type)
	startDate, err := time.Parse(o.DateLayout, o.StartDate)
	if err != nil {
		startDate = time.Now()
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	}

	// Create new watcher
	w := &watcher.Watcher{
		Name:       args[0],
		Type:       eType,
		DatabaseId: o.DatabaseId,
		StartDate:  &startDate,
		WebHook:    o.WebHook,
		Cron:       o.Cron,
		Token:      o.Token,
		Inactive:   o.Inactive,
	}

	// Validate the watcher
	err = w.Validate()
	if err != nil {
		return err
	}

	// Create the watcher
	err = w.Create()
	if err != nil {
		return err
	}

	// Add the newly created watcher to the list of watchers
	watchers.Add(w)

	log.Info().Msg("Watcher created successfully!")

	return nil
}
