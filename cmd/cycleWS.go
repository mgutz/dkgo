package cmd

import (
	"github.com/mgutz/dkgo/pkg/dk"
)

type CycleWSCmd struct {
	Direction string `arg:"" help:"Direction to traverse {${enum}}" enum:"prev,next" default:"next"`
}

func (wc *CycleWSCmd) Run(ctx *Context) error {
	dk.CycleDynamicWS(wc.Direction)
	return nil
}
