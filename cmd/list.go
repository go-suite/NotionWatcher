package cmd

import (
	"fmt"
	nwConfig "github.com/gennesseaux/NotionWatcher/setup/config"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
)

// Options of the root command
type listOptions struct {
	config *nwConfig.NwConfig
}

// root command
func newListCmd() *cobra.Command {

	// Options
	o := listOptions{
		config: nwConfig.Config,
	}

	// Command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list of all watchers currently installed",
		RunE:  o.listCmd,
	}

	return listCmd
}

func (o *listOptions) listCmd(cmd *cobra.Command, args []string) (err error) {

	// Render
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Type", "Database", "Cron", "Hook"})
	for _, w := range watchers.Watchers {
		t.AppendRow(table.Row{w.Name, w.Type, w.DatabaseName, w.Cron, w.WebHook})
	}
	t.SetStyle(table.StyleColoredBright)
	t.Style().Format.Header = text.FormatTitle
	fmt.Println("")
	t.Render()
	fmt.Println("")

	return
}
