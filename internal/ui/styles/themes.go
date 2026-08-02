package styles

import (
	"charm.land/lipgloss/v2"
)

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
