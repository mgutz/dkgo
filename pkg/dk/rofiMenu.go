package dk

import (
	"fmt"
	"strings"

	"github.com/mgutz/dkgo/pkg/util"
)

type RofiMenu struct {
	Clients               []Client
	HasMultipleWorkspaces bool
}

func (m *RofiMenu) Add(client Client) {
	m.Clients = append(m.Clients, client)
}

func (m *RofiMenu) AddNonSelectable(title string) {
	m.Clients = append(m.Clients, Client{WMName: title})
}

func (m *RofiMenu) RofiList() []string {
	var list []string
	for _, client := range m.Clients {
		list = append(list, m.RofiItem(client))
	}
	return list
}

func (m *RofiMenu) RofiItem(client Client) string {
	if client.WMName == "" {
		return "---"
	}

	if client.WMName == "hide" {
		return fmt.Sprintf("%s\x00nonselectable\x1ftrue", client.WMName)
	}

	workspaceInfo := ""
	if m.HasMultipleWorkspaces {
		workspaceInfo = fmt.Sprintf("<span size='xx-small'><tt>%2d</tt></span> ", client.WorkspaceID)
	}

	return fmt.Sprintf("%s<span size='xx-small'><tt>%-9s</tt></span> %s\x00icon\x1f%s", workspaceInfo, strings.ToLower(client.WMClass), client.WMName, client.Icon)
}

func (m *RofiMenu) Run() *Client {
	idx := m.RunRofi()
	if idx == -1 {
		return nil
	}
	return &m.Clients[idx]
}

func (m *RofiMenu) RunRofi() int {
	list := m.RofiList()
	if len(list) == 0 {
		return -1
	}

	lines := util.Clamp(len(list), 3, 16)

	command := fmt.Sprintf("rofi -dmenu -show-icons -markup-rows -format i -theme-str 'listview { lines: %d; }'", lines)
	result, err := util.ExecWithStdin(command, strings.Join(list, "\n"))
	if err != nil {
		return -1
	}

	return util.ParseIntOr(strings.TrimSpace(string(result.Output)), -1)
}
