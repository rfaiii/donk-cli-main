package model

import (
	"image"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/ultraviolet/layout"
	"github.com/richavery/bvr-cli/internal/home"
	"github.com/richavery/bvr-cli/internal/ui/anim"
	"github.com/richavery/bvr-cli/internal/workspace"
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
	cwd := home.Short(m.com.Workspace.WorkingDir())

	accent := t.ThemeColor.Accent
	alt := t.ThemeColor.Alt
	if accent == nil {
		accent = lipgloss.Color("#FF4FA3")
	}
	if alt == nil {
		alt = lipgloss.Color("#B56CFF")
	}

	// Bold, underlined project location rendered in the accent color so it
	// reads clearly against the dark background.
	cwdStyled := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Underline(true).
		Render(cwd)

	// Home buttons: Command palette launcher (top) and File Finder (bottom),
	// each prefixed with a nerd-font glyph from the superfile icon set so
	// they read as actionable buttons. They are stacked vertically with a
	// blank line between them for cushioning, left-aligned one above the
	// other.
	iconColor := t.Header.LogoGradToColor
	terminalIcon := lipgloss.NewStyle().Foreground(iconColor).Render("\ue795") // superfile icon.Terminal
	folderIcon := lipgloss.NewStyle().Foreground(iconColor).Render("\uf07b")   // superfile icon.Directory

	buttonStyle := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(0, 1)

	commandButton := buttonStyle.Render(terminalIcon + " " + "OPEN COMMANDS — ctrl+p")
	folderButton := buttonStyle.Render(folderIcon + " " + "OPEN FILE FINDER — ctrl+shift+f")
	buttons := commandButton + "\n\n" + folderButton

	// Prominent MODEL / PROVIDER line on the homescreen so it's immediately
	// visible. Uses ACCENT for the values and ALT for the labels.
	modelName := ""
	providerName := ""
	if lm := m.selectedLargeModel(); lm != nil {
		modelName = lm.CatwalkCfg.Name
		if pcfg, ok := m.com.Config().Providers.Get(lm.ModelCfg.Provider); ok {
			providerName = pcfg.Name
		}
	}
	modelLine := lipgloss.NewStyle().Foreground(alt).Render("MODEL  ") +
		lipgloss.NewStyle().Foreground(accent).Bold(true).Render(modelName) +
		lipgloss.NewStyle().Foreground(alt).Render("   PROVIDER  ") +
		lipgloss.NewStyle().Foreground(accent).Bold(true).Render(providerName)

	// Click rectangles for the stacked home buttons (Command on top,
	// File Finder below). Both start at the left edge; the File Finder
	// rectangle is pushed down past the Command button plus the gap line.
	btnTop := m.layout.main.Min.Y + 3
	cmdH := lipgloss.Height(commandButton)
	m.commandButtonRect = image.Rect(
		m.layout.main.Min.X, btnTop,
		m.layout.main.Min.X+lipgloss.Width(commandButton), btnTop+cmdH,
	)
	m.finderButtonRect = image.Rect(
		m.layout.main.Min.X, btnTop+cmdH+1,
		m.layout.main.Min.X+lipgloss.Width(folderButton), btnTop+cmdH+1+lipgloss.Height(folderButton),
	)
	parts := []string{cwdStyled, "", buttons, "", modelLine}

	parts = append(parts, "", m.modelInfo(width), "", m.nodeInfo(min(42, width), 3))
	// Idle beaver mascot beneath the status monitors. It walks in step with
	// the banner ticker (m.bannerFrame) so it costs no extra tick loop, and
	// reuses the boot splash frames so it matches the intro.
	parts = append(parts, "", anim.BeaverFrame(m.bannerFrame))
	infoSection := lipgloss.JoinVertical(lipgloss.Left, parts...)

	var remainingHeightArea image.Rectangle
	layout.Vertical(
		layout.Len(lipgloss.Height(infoSection)+1),
		layout.Fill(1),
	).Split(m.layout.main).Assign(new(image.Rectangle), &remainingHeightArea)

	// Left column: NODE, MCP, LSP stacked vertically (top to bottom). The right
	// side holds SKILLS. These are short status monitors, so they share a narrow
	// left column, leaving the bulk of the width for the (often long) skill list.
	leftW := min(30, (width-2)/3)
	rightW := max(1, width-leftW-1)
	sectionH := max(1, remainingHeightArea.Dy())

	nodeSection := m.nodeInfo(leftW, sectionH)
	mcpSection := m.mcpInfo(leftW, sectionH, false)
	lspSection := m.lspInfo(leftW, sectionH, false)
	skillsSection := m.skillsInfo(rightW, sectionH, false)

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, nodeSection, " ", mcpSection, " ", lspSection)
	content := lipgloss.JoinHorizontal(lipgloss.Left, leftColumn, " ", skillsSection)

	return lipgloss.NewStyle().
		Width(width).
		Height(m.layout.main.Dy() - 1).
		PaddingTop(1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, infoSection, "", content),
		)
}
