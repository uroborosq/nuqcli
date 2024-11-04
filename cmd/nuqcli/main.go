package main

import (
	"fmt"

	"github.com/Wifx/gonetworkmanager"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type MainMenuItem struct {
	Label string
}

func (i MainMenuItem) FilterValue() string { return i.Label }

func (i MainMenuItem) Title() string { return i.Label }

func (i MainMenuItem) Description() string { return "" }

type item struct {
	device gonetworkmanager.Device
}

func (i item) Title() string { return string(i.device.GetPath()) }

func (i item) Description() string {
	prop, err := i.device.GetPropertyInterface()
	if err != nil {
		return fmt.Sprintf("error occurred: %s", err.Error())
	}

	return string(prop)
}

func (i item) FilterValue() string {
	return string(i.device.GetPath())
}

type MainMenu struct {
	list list.Model
}

func NewModel() MainMenu {
	items := []list.Item{
		MainMenuItem{Label: "Wi-Fi connections"},
		MainMenuItem{Label: "VPN connections"},
	}
	return MainMenu{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyEnter {
			return NewVPNConnections(), nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenu) View() string {
	return docStyle.Render(m.list.View())
}

type VpnConnectionItem struct {
	vpn gonetworkmanager.VpnConnection
}

func (i VpnConnectionItem) Title() string {
	banner, _ := i.vpn.GetPropertyBanner()
	return banner
}

func (i VpnConnectionItem) Description() string {
	return ""
}

func (i VpnConnectionItem) FilterValue() string {
	return i.Title()
}

type vpnList struct {
	vpnConnections []gonetworkmanager.VpnConnection
	list           list.Model
}

func NewVPNConnections() vpnList {
	return vpnList{list: list.New(nil, list.NewDefaultDelegate(), 0, 0)}
}

func (m vpnList) View() string {
	return docStyle.Render(m.list.View())
}

func (m vpnList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m vpnList) Init() tea.Cmd {
	return nil
}

func main() {
	manager, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		panic(err)
	}

	devices, err := manager.GetAllDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		_, casted := device.(gonetworkmanager.VpnConnection)
		fmt.Println(casted)
		fmt.Printf("%T\n", device)
		fmt.Println(device.GetPropertyInterface())
	}

	// program := tea.NewProgram(NewModel())
	// _, err = program.Run()
	// // fmt.Println(m)
	// if err != nil {
	// 	panic(err)
	// }
}
