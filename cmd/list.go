package cmd

import (
	"encoding/json"
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/modules/config"
	"github.com/gennesseaux/NotionWatcher/modules/watcher"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
)

// Options of the root command
type listOptions struct {
	config *nwConfig.NwConfig

	json bool
}

// list command
func newListCmd() *cobra.Command {

	// Options
	o := listOptions{
		config: nwConfig.Config,
	}

	// Command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list of all watchers currently installed",
		Run:   o.listCmd,
	}

	// Flags
	listCmd.Flags().BoolVarP(&o.json, "json", "j", false, "display the output in json format")

	return listCmd
}

func (o *listOptions) listCmd(_ *cobra.Command, _ []string) {
	if o.json {
		o.outputJSON(watchers.Watchers)
	} else {
		o.outputConsole(watchers.Watchers)
	}
}

// outputConsole outputs the list of watchers in a table
func (o *listOptions) outputConsole(watchers []*watcher.Watcher) {
	// Render
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Type", "Database", "Cron"})
	for _, w := range watchers {
		t.AppendRow(table.Row{w.Name, w.Type, w.DatabaseName, w.Cron})
	}
	t.SetStyle(table.StyleColoredBright)
	t.Style().Format.Header = text.FormatTitle
	fmt.Println("")
	t.Render()
	fmt.Println("")
}

// outputJSON outputs the list of watchers in JSON format
func (o *listOptions) outputJSON(watchers []*watcher.Watcher) {
	j, _ := json.Marshal(watchers)
	fmt.Println(string(j))
}
