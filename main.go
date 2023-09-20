package main

import (
	"github.com/gennesseaux/NotionWatcher/cmd"
	"os"

	// Call implicit init methods
	_ "github.com/gennesseaux/NotionWatcher/modules"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
