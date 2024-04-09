package dk

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/v3/process"

	"github.com/mgutz/dkgo/pkg/desktop"
)

type X11Client struct {
	ID          string
	Commandline string
	PID         int
	WMInstance  string
	WMClass     string
	WMName      string
	WorkspaceID int
	Icon        string
}

type Terminal struct {
	Terminal    string
	TerminalPID int
	Shell       string
	ShellPID    int
	App         string
	AppPID      int
	CommandLine string
}

// RofiWindowList lists all windows in a Rofi menu and swaps the master with the
// selected window. If the selected window is in a different workspace, the
// workspace is focused first.
func RofiWindowList() error {
	status, err := GetStatus()
	if err != nil {
		return err
	}

	wks := status.Global.Focused.Workspace
	dkClients := wks.FocusStack
	focusedClientID := ""
	if len(dkClients) > 0 {
		focusedClientID = dkClients[0].ID
	}

	for _, client := range dkClients {
		if client.Focused {
			focusedClientID = client.ID
			break
		}
	}

	menu := RofiMenu{}
	// put the focused client at the top of the list
	addClients(&menu, dkClients, focusedClientID)
	// add rest of the clients
	for _, wks := range status.Workspaces {
		if wks.Focused {
			continue
		}
		clients := wks.FocusStack
		if len(clients) > 0 {
			menu.HasMultipleWorkspaces = true
			addClients(&menu, clients, "")
		}
	}

	selectedClient := menu.Run()
	if selectedClient == nil {
		return nil
	}

	return SwapMaster(selectedClient.ID)
}

// getCommandline gets the commandline of the process with the given PID.
func getCommandLine(pid int) (string, error) {
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdlineBytes, err := os.ReadFile(cmdlineFile)
	if err != nil {
		return "", fmt.Errorf("failed to read cmdline file: %w", err)
	}

	cmdline := strings.ReplaceAll(string(cmdlineBytes), "\x00", "")
	return cmdline, nil
}

// addClients adds the clients to the Rofi menu.
func addClients(menu *RofiMenu, dkClients []*Client, skipID string) {
	for _, dkClient := range dkClients {
		// do not include the focused client in the list
		if dkClient.ID == skipID {
			continue
		}

		client, err := fromDKClient(dkClient)
		if err != nil || client == nil {
			continue
		}
		menu.Add(*client)
	}
}

// fromDKClient creates a Client from a DK client JSON.
func fromDKClient(dkClient *Client) (*X11Client, error) {
	var icon string
	var err error
	wmInstance := dkClient.Instance
	wmClass := dkClient.Class
	pid := dkClient.PID

	terminal := checkTerminal(pid)
	if terminal != nil {
		icon, err = desktop.GetIcon(terminal.App, terminal.CommandLine)
		if err != nil {
			return nil, fmt.Errorf("getting icon for terminal app %s: %w", terminal.App, err)
		}
	}

	commandline, err := getCommandLine(pid)
	if err != nil {
		return nil, err
	}
	if icon == "" {
		icon = desktop.GetIconOr(wmClass, commandline, "application-x-executable")
	}

	return &X11Client{
		ID:          dkClient.ID,
		WMInstance:  wmInstance,
		WMClass:     wmClass,
		WMName:      dkClient.Title,
		WorkspaceID: dkClient.Workspace,
		PID:         pid,
		Commandline: commandline,
		Icon:        icon,
	}, nil
}

// checkTerminal checks if the process with the given PID is a terminal and
// returns the terminal, shell, and app information.
func checkTerminal(terminalPID int) *Terminal {
	result := &Terminal{}

	terminal, err := process.NewProcess(int32(terminalPID))
	if err != nil {
		return nil
	}

	name, err := terminal.Name()
	if err != nil {
		return nil
	}

	result.Terminal = strings.ToLower(name)
	result.TerminalPID = terminalPID

	if matched, _ := regexp.MatchString("^(st|st-256color|urxvt|kitty|alacritty|xterm|xterm-256colors|wezterm-gui|xfce4-terminal)$", name); !matched {
		return nil
	}

	children, err := terminal.Children()
	if err != nil || len(children) != 1 {
		return nil
	}

	shell := children[0]
	name, err = shell.Name()
	if err != nil {
		return nil
	}
	if matched, _ := regexp.MatchString("zsh|bash|fish", name); !matched {
		return nil
	}

	result.Shell = strings.ToLower(name)
	result.ShellPID = int(shell.Pid)

	children, err = shell.Children()
	if err != nil || len(children) != 1 {
		return nil
	}

	app := children[0]
	name, err = app.Name()
	if err != nil {
		return nil
	}

	// application running in the shell
	result.App = strings.ToLower(name)
	result.AppPID = int(app.Pid)
	commandLine, err := getCommandLine(result.AppPID)
	if err != nil {
		return nil
	}
	result.CommandLine = commandLine

	return result
}
