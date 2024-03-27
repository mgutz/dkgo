package cmd

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/mgutz/dkgo/pkg/util"
)

type Context struct {
	Debug bool
}

var cli struct {
	Debug      bool          `help:"Enable debug mode"`
	Status     StatusCmd     `cmd:"" help:"Show status"`
	SwapMaster SwapCmd       `cmd:"" help:"Swap master on current workspace"`
	CycleWS    CycleWSCmd    `cmd:"" help:"Cycle through dynamic workspace"`
	RofiMaster RofiMasterCmd `cmd:"" help:"Swap master from rofi list, changes workspace as needed"`
}

// Main is the entry point of the application.
func Main() {
	if os.Getenv("DEBUG") == "1" {
		util.CustomizeSlog(util.Untildify("~/.local/state/dkgo.log"))
	}
	kctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := kctx.Run(&Context{Debug: cli.Debug})
	kctx.FatalIfErrorf(err)
}
