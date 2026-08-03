package dialog

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/crush/internal/ui/common"
	"github.com/charmbracelet/crush/internal/ui/styles"
	uv "github.com/charmbracelet/ultraviolet"
)

const ThemesID = "themes"

type ActionSelectTheme struct{ ID string }

type Themes struct {
	com      *common.Common
	selected int
}

func NewThemes(com *common.Common, current string) *Themes {
	selected := 0
	for i, theme := range styles.Themes() {
		if theme.ID == current {
			selected = i
		}
	}
	return &Themes{com: com, selected: selected}
}
func (t *Themes) ID() string { return ThemesID }
func (t *Themes) HandleMsg(msg tea.Msg) Action {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return nil
	}
	themes := styles.Themes()
	switch key.String() {
	case "esc":
		return ActionClose{}
	case "left", "h":
		t.selected = (t.selected + len(themes) - 1) % len(themes)
	case "right", "l":
		t.selected = (t.selected + 1) % len(themes)
	case "enter":
		return ActionSelectTheme{ID: themes[t.selected].ID}
	}
	return nil
}
func (t *Themes) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	themes := styles.Themes()
	lines := []string{"DONK THEMES", "", "←/→ choose  enter apply  esc close", ""}
	for i, theme := range themes {
		marker := "○"
		if i == t.selected {
			marker = "●"
		}
		lines = append(lines, fmt.Sprintf("%s %s", marker, theme.Name))
	}
	view := t.com.Styles.Dialog.View.Width(min(56, max(1, area.Dx()-4))).Render(strings.Join(lines, "\n"))
	DrawCenter(scr, area, view)
	return nil
}
