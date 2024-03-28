package dk

import (
	"fmt"
	"strings"

	"github.com/mgutz/dkgo/pkg/util"
)

type RofiMenu struct {
	Clients               []X11Client
	HasMultipleWorkspaces bool
}

func (m *RofiMenu) Add(client X11Client) {
	m.Clients = append(m.Clients, client)
}

func (m *RofiMenu) AddNonSelectable(title string) {
	m.Clients = append(m.Clients, X11Client{WMName: title})
}

func (m *RofiMenu) menuList() []string {
	var list []string
	for _, client := range m.Clients {
		list = append(list, m.menuItem(client))
	}
	return list
}

func htmlEscape(s string) string {
	return strings.ReplaceAll(s, "<", "&lt;")
}

func (m *RofiMenu) menuItem(client X11Client) string {
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

	return fmt.Sprintf("%s<span size='xx-small'><tt>%-9s</tt></span> %s\x00icon\x1f%s", workspaceInfo, strings.ToLower(client.WMClass), htmlEscape(client.WMName), client.Icon)
}

func (m *RofiMenu) Run() *X11Client {
	idx := m.runRofi()
	if idx == -1 {
		return nil
	}
	return &m.Clients[idx]
}

func (m *RofiMenu) runRofi() int {
	list := m.menuList()
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
