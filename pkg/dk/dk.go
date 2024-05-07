package dk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mgutz/dkgo/pkg/util"
)

// GetStatus returns the status of dk WM.
func GetStatus() (*Status, error) {
	result, err := util.Exec("dkcmd status type=full num=1")
	if err != nil {
		return nil, err
	}

	// remove newlines from output
	output := bytes.Replace(result.Output, []byte("\n"), []byte(""), -1)

	var status Status
	err = json.Unmarshal(output, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// ViewWorkspace switches to the workspace of the client with the given ID.
func ViewWorkspace(clientID string) bool {
	if clientID == "" || !strings.HasPrefix(clientID, "0x") {
		return false
	}

	status, err := GetStatus()
	if err != nil {
		fmt.Println(err)
		return false
	}

	clients := status.Clients
	var wksNumber int
	for _, client := range clients {
		if client.ID == clientID {
			wksNumber = client.Workspace
			break
		}
	}
	if wksNumber == 0 {
		return false
	}

	focusedWorkspaceNumber := status.Global.Focused.Workspace.Number
	if wksNumber != focusedWorkspaceNumber {
		_, err := util.Execf("dkcmd ws %d", wksNumber)
		if err == nil {
			return true
		}
	}
	return false
}

// SwapMaster swaps the master window. v can be a window ID, a tile index, or
// blank, in which case last focused or current focused is used. Master has
// focus after the swap.
func SwapMaster(v string) error {
	// workspaceChanged := ViewWorkspace(v)
	// isID := false
	ViewWorkspace(v)

	status, err := GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	wks := status.Global.Focused.Workspace
	clients := wks.Clients
	// return if no other clients to swap with
	if len(clients) < 2 {
		return nil
	}

	master := clients[0]
	masterID := master.ID
	index := -1

	// v can be a window ID, a tile index, or blank
	if v == "" {
		// Do nothing
	} else if strings.HasPrefix(v, "0x") {
		for i, client := range clients {
			if client.ID == v {
				// isID = true
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

	var target *Client
	if 0 < index && index < len(clients) {
		// caller passed a tile index
		target = clients[index]
	} else if master.Focused && master.Float {
		// allow focused, floating windows to be swapped
		_, err := util.Execf("dkcmd win %s float", masterID)
		if err != nil {
			return err
		}
		return SwapMaster(masterID)
	} else if master.Focused {
		// nothing was passed and master has focus, use last focused window
		target = wks.FocusStack[1]
	} else {
		// nothing was passed, use focused slave
		for _, client := range clients {
			if client.Focused {
				target = client
				break
			}
		}
	}

	if target == nil {
		return nil
	}
	commands := []string{}

	if target.Float {
		util.Execf("dkcmd win %s float", target.ID)
		return SwapMaster(target.ID)
	}

	if target.Focused {
		// ensure current master is recorded in focus history
		commands = append(commands, fmt.Sprintf("dkcmd win %s focus", masterID))
	}
	// swap master with target
	commands = append(commands, fmt.Sprintf("dkcmd win %s swap", target.ID))
	// ensure master has focus
	commands = append(commands, fmt.Sprintf("dkcmd win %s focus", target.ID))

	_, err = util.Bash(strings.Join(commands, " && "))
	return err
}

// CycleDynamicWS cycles through dynamic workspaces like MacOS.
func CycleDynamicWS(direction string) error {
	status, err := GetStatus()
	if err != nil {
		return err
	}

	workspaces := status.Workspaces

	// reversing the order keeps the same logic below for prev, next
	if direction == "prev" {
		for i, j := 0, len(workspaces)-1; i < j; i, j = i+1, j-1 {
			workspaces[i], workspaces[j] = workspaces[j], workspaces[i]
		}
	}

	// determine where to start
	focusedIdx := -1
	for i, ws := range workspaces {
		if ws.Focused {
			focusedIdx = i
			break
		}
	}
	if focusedIdx == -1 {
		return fmt.Errorf("no focused workspace")
	}

	// iterate through workspaces and keep ALL occupied and 1 empty workspace
	emptyFound := false
	for i := 0; i < len(workspaces); i++ {
		idx := (focusedIdx + i) % len(workspaces)
		ws := workspaces[idx]
		if len(ws.Clients) > 0 {
			continue
		}
		if emptyFound {
			ws._skip = true
			continue
		}
		emptyFound = true
	}

	// find the next workspace
	target := -1
	for i := 1; i < len(workspaces); i++ {
		idx := (focusedIdx + i) % len(workspaces)
		ws := workspaces[idx]
		if ws._skip {
			continue
		}
		target = idx
		break
	}
	if target == -1 {
		return fmt.Errorf("no target workspace")
	}

	// view the workspace
	_, err = util.Execf("dkcmd ws %d", workspaces[target].Number)
	return err
}
