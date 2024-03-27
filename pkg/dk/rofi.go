package dk

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mgutz/dkgo/pkg/desktop"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/valyala/fastjson"
)

type Client struct {
	ID          string
	Commandline string
	WMPid       int
	WMInstance  string
	WMClass     string
	WMName      string
	WorkspaceID int
	WMType      string
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

	orderBy := "focus_stack"
	wks := status.Get("global", "focused", "workspace")
	dkClients := wks.GetArray(orderBy)

	focusedClientID := ""
	for _, client := range dkClients {
		if client.GetBool("focused") {
			focusedClientID = string(client.GetStringBytes("id"))
			break
		}
	}

	// Create Rofi menu
	menu := RofiMenu{}
	addClients(&menu, dkClients, focusedClientID)

	for _, wks := range status.GetArray("workspaces") {
		if wks.GetBool("focused") {
			continue
		}
		clients := wks.GetArray(orderBy)
		if len(clients) > 0 {
			menu.HasMultipleWorkspaces = true
			addClients(&menu, clients, "")
		}
	}

	// Run Rofi menu
	selectedClient := menu.Run()
	if selectedClient == nil {
		return nil
	}

	// Swap master with selected client
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
func addClients(menu *RofiMenu, dkClients []*fastjson.Value, skipID string) {
	for _, dkClient := range dkClients {
		// do not include the focused client in the list
		if string(dkClient.GetStringBytes("id")) == skipID {
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
func fromDKClient(dkClient *fastjson.Value) (*Client, error) {
	var icon string
	var err error
	wmInstance := string(dkClient.GetStringBytes("instance"))
	wmClass := string(dkClient.GetStringBytes("class"))
	pid := dkClient.GetInt("pid")

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

	return &Client{
		ID:          string(dkClient.GetStringBytes("id")),
		WMInstance:  wmInstance,
		WMClass:     wmClass,
		WMName:      string(dkClient.GetStringBytes("title")),
		WorkspaceID: dkClient.GetInt("workspace"),
		WMPid:       pid,
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
