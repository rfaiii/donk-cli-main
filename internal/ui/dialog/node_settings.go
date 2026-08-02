package dialog

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/rfaiii/donk-cli-main/internal/node"
	"github.com/rfaiii/donk-cli-main/internal/ui/common"
	"github.com/rfaiii/donk-cli-main/internal/ui/list"
	"github.com/rfaiii/donk-cli-main/internal/ui/styles"
	uv "github.com/charmbracelet/ultraviolet"
)

const (
	nodeSettingsDialogMaxWidth  = 56
	nodeSettingsDialogMaxHeight = 18
)

// nodeSettingsItem is a simple list renderable for the node settings dialog.
type nodeSettingsItem struct {
	device node.Device
	t      *styles.Styles
	ver    uint64
}

func (n nodeSettingsItem) Render(width int) string {
	var icon string
	switch n.device.Status {
	case node.DeviceStatusOnline:
		icon = "●"
	case node.DeviceStatusError:
		icon = "●"
	default:
		icon = "○"
	}
	title := n.device.DisplayName()
	if title == "" {
		title = n.device.ID
	}
	return common.Status(n.t, common.StatusOpts{
		Icon:        icon,
		Title:       title,
		Description: n.device.ConnectionLabel(),
	}, width)
}

func (n nodeSettingsItem) Version() uint64 { return n.ver }
func (n nodeSettingsItem) Finished() bool  { return true }
func (n nodeSettingsItem) Filter() string  { return n.device.DisplayName() }

// NodeSettings represents the node settings dialog.
type NodeSettings struct {
	com        *common.Common
	devices    []node.Device
	selectedID string
	list       *list.List
	loading    bool
	rename     struct {
		id     string
		open   bool
		value  string
		cursor int
	}
}

// NewNodeSettings creates a new node settings dialog.
func NewNodeSettings(com *common.Common) *NodeSettings {
	return &NodeSettings{com: com}
}

// ID implements Dialog.
func (n *NodeSettings) ID() string { return NodeSettingsID }

// Devices updates the visible device list.
func (n *NodeSettings) Devices(devices []node.Device) {
	n.devices = devices
	items := make([]list.Item, 0, len(devices))
	for _, d := range devices {
		items = append(items, nodeSettingsItem{device: d, t: n.com.Styles})
	}
	n.list = list.NewList(items...)
	n.list.SetSize(32, 8)
	if n.selectedID == "" && len(items) > 0 {
		n.selectedID = devices[0].ID
		n.list.SetSelected(0)
	}
}

// HandleMsg implements Dialog.
func (n *NodeSettings) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return ActionClose{}
		case "enter":
			if n.selectedID == "" {
				return nil
			}
			for _, d := range n.devices {
				if d.ID == n.selectedID {
					if strings.EqualFold(d.ConnectionType, "local") {
						return ActionNodeUpdateStatus{DeviceID: d.ID, Status: node.DeviceStatusOnline}
					}
					return ActionNodeRename{DeviceID: d.ID, Current: d.Nickname}
				}
			}
			return nil
		case "up", "ctrl+p":
			if n.list == nil {
				return nil
			}
			if n.list.IsSelectedFirst() {
				n.list.SelectLast()
			} else {
				n.list.SelectPrev()
			}
			n.list.ScrollToSelected()
			syncSelection(n)
		case "down", "ctrl+n":
			if n.list == nil {
				return nil
			}
			if n.list.IsSelectedLast() {
				n.list.SelectFirst()
			} else {
				n.list.SelectNext()
			}
			n.list.ScrollToSelected()
			syncSelection(n)
		case "r":
			if n.selectedID == "" {
				return nil
			}
			for _, d := range n.devices {
				if d.ID == n.selectedID {
					return ActionNodeRename{DeviceID: d.ID, Current: d.Nickname}
				}
			}
			return nil
		}
	case NodeSettingsUpdateMsg:
		n.Devices(msg.Devices)
	}
	return nil
}

func syncSelection(n *NodeSettings) {
	if n.list == nil {
		return
	}
	if selected := n.list.SelectedItem(); selected != nil {
		if item, ok := selected.(nodeSettingsItem); ok {
			n.selectedID = item.device.ID
		}
	}
}

// Cursor returns the desired cursor position.
func (n *NodeSettings) Cursor() *tea.Cursor {
	return &tea.Cursor{}
}

// Draw renders the dialog.
func (n *NodeSettings) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := n.com.Styles
	width := max(0, min(nodeSettingsDialogMaxWidth, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))

	var body string
	if n.loading {
		body = t.Dialog.List.Height(1).Render(t.Resource.AdditionalText.Render("refreshing..."))
	} else if len(n.devices) == 0 {
		body = t.Resource.AdditionalText.Render("None")
	} else if n.list != nil {
		body = t.Dialog.List.Height(n.list.Height()).Render(n.list.Render())
	}

	help := "enter rename • r rename • esc close"
	view := t.Dialog.View.
		Width(width).
		Height(nodeSettingsDialogMaxHeight).
		Render(fmt.Sprintf("%s\n\n%s\n\n%s", t.Dialog.Title.Render("Nodes"), body, t.Dialog.HelpView.Render(help)))

	cur := n.Cursor()
	DrawCenterCursor(scr, area, view, cur)
	return cur
}
