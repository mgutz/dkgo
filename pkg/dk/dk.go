package dk

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mgutz/dkgo/pkg/util"
	"github.com/valyala/fastjson"
)

var (
	FJFalse = fastjson.MustParse("false")
	FJTrue  = fastjson.MustParse("true")
)

// GetStatus returns the status of dk WM.
func GetStatus() (*fastjson.Value, error) {
	result, err := util.Exec("dkcmd status type=json num=1")
	if err != nil {
		return nil, err
	}
	return fastjson.ParseBytes(result.Output)
}

// ViewWorkspace switches to the workspace of the client with the given ID.
func ViewWorkspace(clientID string) {
	if clientID == "" || !strings.HasPrefix(clientID, "0x") {
		return
	}

	status, err := GetStatus()
	if err != nil {
		log.Println(err)
		return
	}

	clients := status.GetArray("clients")
	var wksNumber int
	for _, client := range clients {
		if string(client.GetStringBytes("id")) == clientID {
			wksNumber = client.GetInt("workspace")
			break
		}
	}
	if wksNumber == 0 {
		return
	}

	focusedWorkspaceNumber := status.GetInt("global", "focused", "workspace", "number")
	if wksNumber != focusedWorkspaceNumber {
		util.Execf("dkcmd ws %d", wksNumber)
	}
}

// SwapMaster swaps the master window. v can be a window ID, a tile index, or
// blank, in which case last focused or current focused is used. Master has
// focus after the swap.
func SwapMaster(v string) error {
	ViewWorkspace(v)

	status, err := GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	wks := status.Get("global", "focused", "workspace")
	clients := wks.GetArray("clients")
	// return if no other clients to swap with
	if len(clients) < 2 {
		return nil
	}

	master := clients[0]
	masterID := string(master.GetStringBytes("id"))
	index := -1

	// v can be a window ID, a tile index, or blank
	if v == "" {
		// Do nothing
	} else if strings.HasPrefix(v, "0x") {
		for i, client := range clients {
			if string(client.GetStringBytes("id")) == v {
				index = i
				break
			}
		}
	} else {
		index, err = strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid index argument: %w", err)
		}
	}

	var target *fastjson.Value
	if 0 < index && index < len(clients) {
		// caller passed a tile index
		target = clients[index]
	} else if master.GetBool("focused") && master.GetBool("float") {
		// allow focused, floating windows to be swapped
		_, err := util.Execf("dkcmd win %s float", masterID)
		if err != nil {
			return err
		}
		return SwapMaster(masterID)
	} else if master.GetBool("focused") {
		// nothing was passed and master has focus, use last focused window
		target = wks.Get("focus_stack", "1")
	} else {
		// nothing was passed, use focused slave
		for _, client := range clients {
			if client.GetBool("focused") {
				target = client
				break
			}
		}
	}

	if target == nil {
		return nil
	}
	targetID := string(target.GetStringBytes("id"))
	commands := []string{}

	if target.GetBool("float") {
		util.Execf("dkcmd win %s float", targetID)
		return SwapMaster(targetID)
	}

	if target.GetBool("focused") {
		// ensure current master is recorded in focus history
		commands = append(commands, fmt.Sprintf("dkcmd win %s focus", masterID))
	}
	// swap master with target
	commands = append(commands, fmt.Sprintf("dkcmd win %s swap", targetID))
	// ensure master has focus
	commands = append(commands, fmt.Sprintf("dkcmd win %s focus", targetID))

	_, err = util.Bash(strings.Join(commands, " && "))
	return err
}

// CycleDynamicWS cycles through dynamic workspaces like MacOS.
func CycleDynamicWS(direction string) error {
	status, err := GetStatus()
	if err != nil {
		return err
	}

	// iterate through workspaces and keep ALL occupied and 1 empty workspace
	workspaces := status.GetArray("workspaces")
	emptyFound := false
	for _, ws := range workspaces {
		if len(ws.GetArray("clients")) > 0 {
			continue
		}
		if emptyFound {
			ws.Set("SKIP", FJTrue)
			continue
		}
		emptyFound = true
	}

	// reversing the order keeps the same logic below for prev, next
	if direction == "prev" {
		for i, j := 0, len(workspaces)-1; i < j; i, j = i+1, j-1 {
			workspaces[i], workspaces[j] = workspaces[j], workspaces[i]
		}
	}

	// determine where to start
	focusedIdx := -1
	for i, ws := range workspaces {
		if ws.GetBool("focused") {
			focusedIdx = i
			break
		}
	}
	if focusedIdx == -1 {
		return fmt.Errorf("no focused workspace")
	}

	// find the next workspace to switch to
	target := -1
	for i := 1; i < len(workspaces); i++ {
		idx := (focusedIdx + i) % len(workspaces)
		ws := workspaces[idx]
		if ws.GetBool("SKIP") {
			continue
		}
		target = idx
		break
	}
	if target == -1 {
		return fmt.Errorf("no target workspace")
	}

	// view the workspace
	_, err = util.Execf("dkcmd ws %d", workspaces[target].GetInt("number"))
	return err
}
