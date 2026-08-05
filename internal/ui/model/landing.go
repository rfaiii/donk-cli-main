package model

import (
	"image"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/ultraviolet/layout"
	"github.com/richavery/donk-cli/internal/ui/common"
	"github.com/richavery/donk-cli/internal/workspace"
)

// selectedLargeModel returns the currently selected large language model from
// the agent coordinator, if one exists.
func (m *UI) selectedLargeModel() *workspace.AgentModel {
	if m.com.Workspace.AgentIsReady() {
		model := m.com.Workspace.AgentModel()
		return &model
	}
	return nil
}

// landingView renders the landing page view showing the current working
// directory, model information, and LSP/MCP status in a two-column layout.
func (m *UI) landingView() string {
	t := m.com.Styles
	width := m.layout.main.Dx()
	cwd := common.PrettyPath(t, m.com.Workspace.WorkingDir(), width)

	buttonText := "📁  OPEN FILE FINDER  (ctrl+shift+f)"
	button := lipgloss.NewStyle().Foreground(lipgloss.Color("#3BF66B")).Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3BF66B")).Padding(0, 1).Render(buttonText)
	parts := []string{cwd, "", button}

	m.finderButtonRect = image.Rect(m.layout.main.Min.X, m.layout.main.Min.Y+3, m.layout.main.Min.X+lipgloss.Width(button), m.layout.main.Min.Y+3+lipgloss.Height(button))
	parts = append(parts, "", m.modelInfo(width), "", m.nodeInfo(min(42, width), 3))
	infoSection := lipgloss.JoinVertical(lipgloss.Left, parts...)

	var remainingHeightArea image.Rectangle
	layout.Vertical(
		layout.Len(lipgloss.Height(infoSection)+1),
		layout.Fill(1),
	).Split(m.layout.main).Assign(new(image.Rectangle), &remainingHeightArea)

	mcpLspSectionWidth := min(30, (width-2)/3)

	lspSection := m.lspInfo(mcpLspSectionWidth, max(1, remainingHeightArea.Dy()), false)
	mcpSection := m.mcpInfo(mcpLspSectionWidth, max(1, remainingHeightArea.Dy()), false)
	skillsSection := m.skillsInfo(mcpLspSectionWidth, max(1, remainingHeightArea.Dy()), false)
	nodeSection := m.nodeInfo(mcpLspSectionWidth, max(1, remainingHeightArea.Dy()))

	content := lipgloss.JoinHorizontal(lipgloss.Left, lspSection, " ", mcpSection, " ", skillsSection, " ", nodeSection)

	return lipgloss.NewStyle().
		Width(width).
		Height(m.layout.main.Dy() - 1).
		PaddingTop(1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, infoSection, "", content),
		)
}
