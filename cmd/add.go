package cmd

import (
	"errors"
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/event"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	log "github.com/go-mods/zerolog-quick"
	"github.com/spf13/cobra"
	"time"
)

type addOptions struct {
	config *nwConfig.NwConfig

	Name        string
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

func newAddCmd() *cobra.Command {

	// Options
	o := addOptions{
		config: nwConfig.Config,
	}

	// Command
	addCmd := &cobra.Command{
		Use:     "add [Name]",
		Aliases: []string{"a", "create", "c", "new"},
		Short:   "Add a new watcher",
		Args:    cobra.ExactArgs(1),
		RunE:    o.addCmd,
	}

	// Disable usage print on error
	addCmd.SilenceUsage = true

	// Command flags
	addCmd.Flags().StringVarP(&o.Type, "type", "", "", "type of event. Can be pageAddedToDatabase or pageUpdatedInDatabase")
	addCmd.Flags().StringVarP(&o.DatabaseId, "database-id", "", "", "id of the database")
	addCmd.Flags().StringVarP(&o.StartDate, "start-date", "", "", "minimum date to watch")
	addCmd.Flags().StringVarP(&o.DateLayout, "date-layout", "", "2006-01-02", "layout to use when parsing the start date (default Y-m-d)")
	addCmd.Flags().StringVarP(&o.WebHook, "hook", "", "", "WebHook to call on update")
	addCmd.Flags().StringVarP(&o.WebHookTest, "hook-test", "", "", "WebHook to call for testing")
	addCmd.Flags().StringVarP(&o.Cron, "cron", "", "", "cron")
	addCmd.Flags().StringVarP(&o.Token, "token", "", "", "Notion token")
	addCmd.Flags().BoolVarP(&o.Inactive, "inactive", "", false, "Inactive watcher")

	// Command required flags
	_ = addCmd.MarkFlagRequired("type")
	_ = addCmd.MarkFlagRequired("database-id")
	_ = addCmd.MarkFlagRequired("hook")
	_ = addCmd.MarkFlagRequired("token")

	return addCmd
}

func (o *addOptions) addCmd(_ *cobra.Command, args []string) (err error) {

	// Check if a watcher with the same name exist
	if watchers.Get(args[0]) != nil {
		return errors.New("a watcher with the same name already exist")
	}

	// Parse the type
	eType, _ := event.ParseType(o.Type)

	// Parse the start date
	startDate, err := time.Parse(o.DateLayout, o.StartDate)
	if err != nil {
		startDate = time.Now()
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	}

	// Create new watcher
	w := &watcher.Watcher{
		Name:        args[0],
		Type:        eType,
		DatabaseId:  o.DatabaseId,
		StartDate:   &startDate,
		WebHook:     o.WebHook,
		WebHookTest: o.WebHookTest,
		Cron:        o.Cron,
		Token:       o.Token,
		Inactive:    o.Inactive,
	}

	// Validate the watcher
	err = w.Validate()
	if err != nil {
		return err
	}

	// Create the watcher and return an error if it fails
	err = w.Create()
	if err != nil {
		return fmt.Errorf("failed to create the watcher: %s", err.Error())
	}

	// Add the newly created watcher to the list of watchers
	watchers.Add(w)

	log.Info().Msg("Watcher created successfully!")

	return nil
}
