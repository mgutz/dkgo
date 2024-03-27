package cmd

import (
	"github.com/mgutz/dkgo/pkg/dk"
)

type SwapCmd struct {
	AddressOrTileIndex string `arg:"" default:""`
}

func (wc *SwapCmd) Run(ctx *Context) error {
	dk.SwapMaster(wc.AddressOrTileIndex)
	return nil
}
