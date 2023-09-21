package cmd

import (
	"encoding/json"
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	"github.com/spf13/cobra"
)

// Options of the info command
type infoOptions struct {
	config *nwConfig.NwConfig

	json bool
}

// info command
func newInfoCmd() *cobra.Command {

	// Options
	o := infoOptions{
		config: nwConfig.Config,
	}

	// Command
	infoCmd := &cobra.Command{
		Use:     "info [name]",
		Aliases: []string{"i"},
		Short:   "display information about the watcher",
		RunE:    o.infoCmd,
	}

	// Flags
	infoCmd.Flags().BoolVarP(&o.json, "json", "j", false, "display the output in json format")

	return infoCmd
}

func (o *infoOptions) infoCmd(_ *cobra.Command, args []string) (err error) {

	// Get the watcher or return an error
	w := watchers.Get(args[0])
	if w == nil {
		return fmt.Errorf("the watcher %s does not exist", args[0])
	}

	if o.json {
		o.outputJSON(w)
	} else {
		o.outputConsole(w)
	}

	return nil
}

// outputConsole outputs the watcher information in the console
func (o *infoOptions) outputConsole(w *watcher.Watcher) {
	fmt.Printf("\tName: %s\n", w.Name)
	fmt.Printf("\tType: %s\n", w.Type)
	fmt.Printf("\tDatabase: %s\n", w.DatabaseName)
	fmt.Printf("\tStartDate: %s\n", w.StartDate)
	fmt.Printf("\tCron: %s\n", w.Cron)
	fmt.Printf("\tWebHook: %s\n", w.WebHook)
	fmt.Printf("\tWebHookTest: %s\n", w.WebHookTest)
	fmt.Printf("\tActive: %t\n", !w.Inactive)
}

// outputJSON outputs the watcher information in JSON format
func (o *infoOptions) outputJSON(w *watcher.Watcher) {
	j, _ := json.Marshal(w)
	fmt.Println(string(j))
}
