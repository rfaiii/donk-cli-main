// Package logo renders the DONK-CLI wordmark.
package logo

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/richavery/donk-cli/internal/ui/styles"
)

// letterform represents a letterform. It can be stretched horizontally by
// a given amount via the boolean argument.
type letterform func(bool) string

const diag = `╱`

var bannerGlitchRunes = []rune("01ABCDEF~!@#$%^&*+=-/\\|")

var donkCLIASCII = map[rune][]string{
	'D': {"████ ", "█   █", "█   █", "█   █", "████ "},
	'O': {" ███ ", "█   █", "█   █", "█   █", " ███ "},
	'N': {"█   █", "██  █", "█ █ █", "█  ██", "█   █"},
	'K': {"█  █ ", "█ █  ", "██   ", "█ █  ", "█  █ "},
	'-': {"     ", "     ", "█████", "     ", "     "},
	'C': {" ████", "█    ", "█    ", "█    ", " ████"},
	'L': {"█    ", "█    ", "█    ", "█    ", "█████"},
	'I': {"█████", "  █  ", "  █  ", "  █  ", "█████"},
}

// Wordmark is the application name shown in the terminal UI. Keep this as the
// single source of truth; use scripts/set-wordmark.sh to change it safely.
const Wordmark = "DONK"

// Opts are the options for rendering the DONK-CLI title art.
type Opts struct {
	FieldColor   color.Color // diagonal lines
	TitleColorA  color.Color // left gradient ramp point
	TitleColorB  color.Color // right gradient ramp point
	CharmColor   color.Color // Legacy metadata color retained for API compatibility
	VersionColor color.Color // version text color
	Width        int         // width of the rendered logo, used for truncation
	Hyper        bool        // retained for compatibility; it does not alter the DONK wordmark

	// When true, stretch a random letterform on each render. Has no effect in
	// compact mode. Mainly for testing. In production you will want to cache
	// the stretched letterform to keep the logo from jittering on resize.
	Unstable bool

	// When true, render the title as plain text instead of
	// block characters. Useful when the terminal font does not support the
	// block characters used in the stylized letterforms.
	Text bool

	// Animated enables the glitch strip beneath the wide wordmark.
	Animated bool
	// Frame selects the current glitch frame.
	Frame int
	// Animation is the rendered prompt-style scramble for the banner strip.
	Animation string
}

// donkCLIWordmark is the compact block-letter suffix used by the Donk brand.
// Keeping it in the logo package makes the branding independent from terminal
// fonts while preserving the existing DONK letterform treatment.
func donkCLIWordmark(stretch bool) string {
	cli := "█▀▀▀  █  █\n█     █  █\n▀▀▀▀  ▀  ▀"
	if stretch {
		cli = strings.Replace(cli, "█▀▀▀", "█▀▀▀▀▀", 1)
	}
	return cli
}

// Render renders the DONK-CLI logo. Set the argument to true to render the narrow
// version, intended for use in a sidebar.
//
// The compact argument determines whether it renders compact for the sidebar
// or wider for the main pane.
func Render(base lipgloss.Style, version string, compact bool, o Opts) string {
	fg := func(c color.Color, s string) string {
		return lipgloss.NewStyle().Foreground(c).Render(s)
	}

	// Title.
	var donk string
	var donkWidth int
	if o.Text {
		if compact {
			donk = styles.ApplyForegroundGrad(base, Wordmark, o.TitleColorA, o.TitleColorB)
		} else {
			donk = renderASCIIWordmark(base, o.TitleColorA, o.TitleColorB)
		}
		donkWidth = lipgloss.Width(donk)
	} else {
		const spacing = 1
		donkLetterforms := []letterform{
			LetterD,
			LetterO,
			LetterN,
			LetterK,
		}

		stretchIndex := -1 // -1 means no stretching.
		if !compact && !o.Unstable {
			// Always stretch the same letterform, which is picked once at random.
			stretchIndex = cachedRandN(len(donkLetterforms))
		} else if !compact && o.Unstable {
			// Stretch a random letterform on every render.
			stretchIndex = rand.IntN(len(donkLetterforms))
		}
		donk = renderWord(spacing, stretchIndex, donkLetterforms...)
		cli := donkCLIWordmark(false)
		donk = lipgloss.JoinHorizontal(lipgloss.Top, donk, "-", cli)
		// Hyper is a provider, not a separate DONK product name. The UI keeps
		// one stable wordmark regardless of the selected provider.
		donkWidth = lipgloss.Width(donk)
		b := new(strings.Builder)
		for r := range strings.SplitSeq(donk, "\n") {
			fmt.Fprintln(b, styles.ApplyForegroundGrad(base, r, o.TitleColorA, o.TitleColorB))
		}
		donk = b.String()
	}

	// The Donk-CLI wordmark stands alone; no upstream metadata row.
	_ = version

	// Narrow version.
	if compact {
		field := fg(o.FieldColor, strings.Repeat(diag, donkWidth))
		return strings.Join([]string{field, field, donk, field, ""}, "\n")
	}

	fieldHeight := lipgloss.Height(donk)

	// Left field.
	const leftWidth = 6
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diag, leftWidth))
	leftField := new(strings.Builder)
	for range fieldHeight {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field.
	rightWidth := max(15, o.Width-donkWidth-leftWidth-2) // 2 for the gap.
	const stepDownAt = 0
	rightField := new(strings.Builder)
	for i := range fieldHeight {
		width := rightWidth
		if i >= stepDownAt {
			width = rightWidth - (i - stepDownAt)
		}
		fmt.Fprint(rightField, fg(o.FieldColor, strings.Repeat(diag, width)), "\n")
	}
	leftFieldText := strings.TrimSuffix(leftField.String(), "\n")
	rightFieldText := strings.TrimSuffix(rightField.String(), "\n")

	// Return the wide version.
	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftFieldText, hGap, donk, hGap, rightFieldText)
	if o.Animated {
		glitch := o.Animation
		if glitch == "" {
			glitch = glitchPattern(max(o.Width, lipgloss.Width(logo)), o.Frame)
		}
		glitch = ansi.Truncate(glitch, max(o.Width, lipgloss.Width(logo)), "")
		logo = lipgloss.JoinVertical(lipgloss.Left, logo, glitch)
	}
	if o.Width > 0 {
		// Truncate the logo to the specified width.
		lines := strings.Split(logo, "\n")
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, o.Width, "")
		}
		logo = strings.Join(lines, "\n")
	}
	return logo
}

func glitchPattern(width, frame int) string {
	if width < 1 {
		return ""
	}
	var b strings.Builder
	for i := range width {
		// Characters scroll steadily, while frequent positions rapidly mutate
		// to create a visible terminal glitch/code effect.
		index := (i + frame) % len(bannerGlitchRunes)
		if (i+frame*3)%7 < 2 {
			index = (i*i + frame*5 + i*frame) % len(bannerGlitchRunes)
		}
		b.WriteRune(bannerGlitchRunes[index])
	}
	return b.String()
}

func renderASCIIWordmark(base lipgloss.Style, colorA, colorB color.Color) string {
	rows := make([]string, 5)
	for _, r := range Wordmark {
		glyph, ok := donkCLIASCII[r]
		if !ok {
			continue
		}
		for row := range glyph {
			rows[row] += glyph[row] + " "
		}
	}
	var b strings.Builder
	for i, row := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(styles.ApplyForegroundGrad(base, row, colorA, colorB))
	}
	return b.String()
}

// SmallRender renders a smaller version of the DONK-CLI logo, suitable for
// smaller windows or sidebar usage.
func SmallRender(t *styles.Styles, width int, o Opts) string {
	title := styles.ApplyBoldForegroundGrad(t.Logo.GradCanvas, Wordmark, t.Logo.SmallGradFromColor, t.Logo.SmallGradToColor)
	remainingWidth := width - lipgloss.Width(title) - 1
	if remainingWidth > 0 {
		lines := strings.Repeat("╱", remainingWidth)
		title = fmt.Sprintf("%s %s", title, t.Logo.SmallDiagonals.Render(lines))
	}
	return title
}
