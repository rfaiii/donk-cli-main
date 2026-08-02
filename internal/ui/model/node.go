package model

import (
	"slices"

	"charm.land/lipgloss/v2"
	"github.com/rfaiii/donk-cli-main/internal/node"
	"github.com/rfaiii/donk-cli-main/internal/ui/common"
)

// nodeInfo renders the Node connection status section showing connected devices
// and their connection status with colored status dots.
func (m *UI) nodeInfo(width, maxItems int, isSection bool) string {
	t := m.com.Styles
	devices := node.Devices()
	if len(devices) == 0 {
		return common.Section(t, t.Resource.Heading.Render("Node")+"\n"+t.Resource.AdditionalText.Render("No devices"), width)
	}

	title := t.Resource.Heading.Render("Node")
	var items []string
	for _, d := range devices {
		var icon, color string
		switch d.Status {
		case node.DeviceStatusOnline:
			icon = "●"
			color = "#3BF66B"
		case node.DeviceStatusError:
			icon = "●"
			color = "#FF4444"
		default:
			icon = "○"
			color = "#888888"
		}

		name := d.DisplayName()
		if name == "" {
			name = d.ID
		}

		itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		label := itemStyle.Render(icon + " " + name)
		desc := t.Resource.AdditionalText.Render(d.ConnectionLabel())
		items = append(items, label+" "+desc)
	}

	// Limit displayed items to available height
	if maxItems > 0 && len(items) > maxItems {
		items = slices.Delete(items, maxItems-1, len(items))
		items = append(items, t.Resource.AdditionalText.Render("..."))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, items...)
	return common.Section(t, title+"\n"+body, width)
}
