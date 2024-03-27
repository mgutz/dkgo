package cmd

import (
	"github.com/mgutz/dkgo/pkg/dk"
)

// RofiMasterCmd swaps master from rofi list, changes workspace as needed.
type RofiMasterCmd struct{}

func (wc *RofiMasterCmd) Run(ctx *Context) error {
	return dk.RofiWindowList()
}
