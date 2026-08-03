package styles

import (
	"charm.land/lipgloss/v2"
	"image/color"
)

const DefaultTheme = "rich-aizen-green"

type ThemeDefinition struct {
	ID, Name           string
	Primary, Secondary color.Color
}

var themeDefinitions = []ThemeDefinition{
	{DefaultTheme, "Rich Aizen Green", lipgloss.Color("#3BF66B"), lipgloss.Color("#3BF66B")},
	{"crazy-jeff-pink", "Crazy Jeff Pink", lipgloss.Color("#FF4FA3"), lipgloss.Color("#FF86C8")},
	{"kobe-yang-purple", "Kobe Yang Purple", lipgloss.Color("#B56CFF"), lipgloss.Color("#E0B3FF")},
	{"steve-dabeav-blue", "Steve DaBeav Blue", lipgloss.Color("#5CC8FF"), lipgloss.Color("#A9E7FF")},
	{"jenny-ann-orange", "Jenny Ann Orange", lipgloss.Color("#FF8A3D"), lipgloss.Color("#FFC078")},
	{"felix-tornado-white", "Felix Tornado White", lipgloss.Color("#FFFFFF"), lipgloss.Color("#E8EEF2")},
	{"luis-mellow-yellow", "Luis Mellow Yellow", lipgloss.Color("#D6C84A"), lipgloss.Color("#F2E98B")},
}

func Themes() []ThemeDefinition { return append([]ThemeDefinition(nil), themeDefinitions...) }
func ThemeByID(id string) ThemeDefinition {
	for _, theme := range themeDefinitions {
		if theme.ID == id {
			return theme
		}
	}
	return themeDefinitions[0]
}
func ThemeForName(id string) Styles {
	theme := ThemeByID(id)
	s := DarkDonkTheme()
	return applyThemeAccent(s, theme.Primary, theme.Secondary)
}

func applyThemeAccent(s Styles, primary, secondary color.Color) Styles {
	s.Header.LogoGradFromColor, s.Header.LogoGradToColor = secondary, primary
	s.Header.HypercreditIcon = s.Header.HypercreditIcon.Foreground(secondary)
	s.Header.Percentage = s.Header.Percentage.Foreground(primary)
	s.Editor.PromptNormalFocused = s.Editor.PromptNormalFocused.Foreground(primary)
	s.Editor.PromptBeastmodeIconFocused = s.Editor.PromptBeastmodeIconFocused.Foreground(primary)
	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.Foreground(primary)
	s.ModelInfo.HypercreditIcon = s.ModelInfo.HypercreditIcon.Foreground(secondary)
	s.Resource.OnlineIcon = s.Resource.OnlineIcon.Foreground(primary)
	s.Dialog.Title = s.Dialog.Title.Foreground(primary)
	s.Dialog.TitleText = s.Dialog.TitleText.Foreground(primary)
	s.Dialog.SelectedItem = s.Dialog.SelectedItem.BorderForeground(primary)
	s.Messages.UserFocused = s.Messages.UserFocused.BorderForeground(primary)
	s.Messages.ShellBarFocused = s.Messages.ShellBarFocused.BorderForeground(primary)
	s.Logo.TitleColorA, s.Logo.TitleColorB = secondary, primary
	s.Logo.FieldColor = primary
	return s
}

// ThemeKeyForProvider returns a stable identifier for the theme
// associated with the given provider ID. Providers that share a theme
// yield the same key, so callers can cheaply detect when switching
// providers would not actually change the active theme and skip the
// expensive style rebuild. This is the single source of truth for the
// provider-to-theme mapping; [ThemeForProvider] builds on it.
func ThemeKeyForProvider(providerID string) string {
	switch providerID {
	case "hyper":
		return "dark-donk-theme"
	default:
		return "dark-donk-theme"
	}
}

// ThemeForProvider returns the Dark Donk theme for every provider. Provider
// identity is retained for compatibility, but DONK has one visual identity.
func ThemeForProvider(providerID string) Styles {
	return DarkDonkTheme()
}

// DarkDonkTheme returns DONK's dark green theme. The palette is shared with
// the cross-language references in R&D/DONK-DARK-COLOR.
func DarkDonkTheme() Styles {
	s := quickStyle(quickStyleOpts{
		primary:   lipgloss.Color("#3BF66B"),
		secondary: lipgloss.Color("#3BF66B"),
		accent:    lipgloss.Color("#FFFFFF"),
		keyword:   lipgloss.Color("#3BF66B"),

		fgBase:       lipgloss.Color("#FFFFFF"),
		fgMoreSubtle: lipgloss.Color("#8C9691"),
		fgSubtle:     lipgloss.Color("#8C9691"),
		fgMostSubtle: lipgloss.Color("#59645E"),

		onPrimary: lipgloss.Color("#0C0E0D"),

		bgBase:         lipgloss.Color("#0C0E0D"),
		bgLeastVisible: lipgloss.Color("#0B0D0C"),
		bgLessVisible:  lipgloss.Color("#141716"),
		bgMostVisible:  lipgloss.Color("#222825"),

		separator: lipgloss.Color("#222825"),

		destructive:       lipgloss.Color("#FF5F56"),
		error:             lipgloss.Color("#FF5F56"),
		warningSubtle:     lipgloss.Color("#F2C14E"),
		warning:           lipgloss.Color("#F2C14E"),
		denied:            lipgloss.Color("#D98B5F"),
		busy:              lipgloss.Color("#3BF66B"),
		info:              lipgloss.Color("#8CDED0"),
		infoMoreSubtle:    lipgloss.Color("#5E9B8E"),
		infoMostSubtle:    lipgloss.Color("#37635A"),
		success:           lipgloss.Color("#3BF66B"),
		successMoreSubtle: lipgloss.Color("#3BF66B"),
		successMostSubtle: lipgloss.Color("#1B3B2B"),

		// ANSI 16-color palette for remapping raw terminal output
		// (e.g. bang-mode shell commands) onto legible Charmtone colors.
		ansiBlack:   lipgloss.Color("#0C0E0D"),
		ansiRed:     lipgloss.Color("#FF5F56"),
		ansiGreen:   lipgloss.Color("#3BF66B"),
		ansiYellow:  lipgloss.Color("#F2C14E"),
		ansiBlue:    lipgloss.Color("#8CDED0"),
		ansiMagenta: lipgloss.Color("#C69CFF"),
		ansiCyan:    lipgloss.Color("#8CDED0"),
		ansiWhite:   lipgloss.Color("#FFFFFF"),

		ansiBrightBlack:   lipgloss.Color("#59645E"),
		ansiBrightRed:     lipgloss.Color("#FF7B72"),
		ansiBrightGreen:   lipgloss.Color("#6BFF91"),
		ansiBrightYellow:  lipgloss.Color("#FFD866"),
		ansiBrightBlue:    lipgloss.Color("#B0FFF0"),
		ansiBrightMagenta: lipgloss.Color("#E0C2FF"),
		ansiBrightCyan:    lipgloss.Color("#B0FFF0"),
		ansiBrightWhite:   lipgloss.Color("#FFFFFF"),
	})

	// Bang ! prompt overrides use the Dark Donk green/white pairing.
	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#3BF66B"))
	s.Editor.PromptBangDotsFocused = s.Editor.PromptBangDotsFocused.
		Foreground(lipgloss.Color("#3BF66B"))
	s.Editor.PromptBangDotsBlurred = s.Editor.PromptBangDotsBlurred.
		Foreground(lipgloss.Color("#59645E"))

	// Shell bar/prompt overrides use the same green accent.
	s.Messages.ShellBarFocused = s.Messages.ShellBarFocused.
		BorderForeground(lipgloss.Color("#3BF66B"))
	s.Messages.ShellBarBlurred = s.Messages.ShellBarBlurred.
		BorderForeground(lipgloss.Color("#222825"))
	s.Messages.ShellPrompt = s.Messages.ShellPrompt.
		Foreground(lipgloss.Color("#3BF66B"))
	s.Messages.ShellPromptBlurred = s.Messages.ShellPromptBlurred.
		Foreground(lipgloss.Color("#8C9691"))

	return s
}

// CharmtonePantera is retained as a source-compatible alias for existing UI
// tests and callers. New code should use DarkDonkTheme.
func CharmtonePantera() Styles { return DarkDonkTheme() }

// HypercrushObsidiana returns the Hypercrush dark theme.
func HypercrushObsidiana() Styles {
	return DarkDonkTheme()
}
