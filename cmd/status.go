package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mgutz/dkgo/pkg/dk"
)

// Status pretty prints dkcmd status output.
type StatusCmd struct{}

func (sc *StatusCmd) Run(ctx *Context) error {
	status, err := dk.GetStatus()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
