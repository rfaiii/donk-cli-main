// Package anim provides Charm Bubbles–style spinner presets ported from
// the Rust donk-anim crate.
package anim

import "time"

// SpinnerKind identifies a spinner preset.
type SpinnerKind string

const (
	SpinnerKindLine     SpinnerKind = "line"
	SpinnerKindDot      SpinnerKind = "dot"
	SpinnerKindMiniDot  SpinnerKind = "minidot"
	SpinnerKindJump     SpinnerKind = "jump"
	SpinnerKindPulse    SpinnerKind = "pulse"
	SpinnerKindPoints   SpinnerKind = "points"
	SpinnerKindMoon     SpinnerKind = "moon"
)

// SpinnerKinds is the ordered list of available spinner presets.
var SpinnerKinds = []SpinnerKind{
	SpinnerKindLine,
	SpinnerKindDot,
	SpinnerKindMiniDot,
	SpinnerKindJump,
	SpinnerKindPulse,
	SpinnerKindPoints,
	SpinnerKindMoon,
}

// spinnerFrames maps each kind to its frame set.
var spinnerFrames = map[SpinnerKind][]string{
	SpinnerKindLine:    {"|", "/", "-", "\\"},
	SpinnerKindDot:     {"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
	SpinnerKindMiniDot: {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	SpinnerKindJump:    {"⢄", "⢂", "⢁", "⡁", "⡈", "⡐", "⡠"},
	SpinnerKindPulse:   {"▏", "▎", "▍", "▌", "▋", "▊", "▉", "▊", "▋", "▌", "▍", "▎"},
	SpinnerKindPoints:  {"∙∙∙", "●∙∙", "∙●∙", "∙∙●"},
	SpinnerKindMoon:    {"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
}

// spinnerInterval maps each kind to its tick interval.
var spinnerInterval = map[SpinnerKind]time.Duration{
	SpinnerKindLine:    100 * time.Millisecond,
	SpinnerKindDot:     80 * time.Millisecond,
	SpinnerKindMiniDot: 80 * time.Millisecond,
	SpinnerKindJump:    100 * time.Millisecond,
	SpinnerKindPulse:   120 * time.Millisecond,
	SpinnerKindPoints:  160 * time.Millisecond,
	SpinnerKindMoon:    200 * time.Millisecond,
}

// Spinner is a Charm Bubbles–style spinner frame player.
type Spinner struct {
	kind    SpinnerKind
	index   int
	frames  []string
	last    time.Time
	interval time.Duration
}

// NewSpinner creates a spinner with the given preset.
func NewSpinner(kind SpinnerKind) *Spinner {
	frames := spinnerFrames[kind]
	if len(frames) == 0 {
		frames = spinnerFrames[SpinnerKindDot]
		kind = SpinnerKindDot
	}
	return &Spinner{
		kind:     kind,
		index:    0,
		frames:   frames,
		last:     time.Now(),
		interval: spinnerInterval[kind],
	}
}

// Kind returns the active spinner preset.
func (s *Spinner) Kind() SpinnerKind { return s.kind }

// SetKind changes the spinner preset and resets position.
func (s *Spinner) SetKind(kind SpinnerKind) {
	frames := spinnerFrames[kind]
	if len(frames) == 0 {
		return
	}
	s.kind = kind
	s.frames = frames
	s.index = 0
	s.last = time.Now()
	s.interval = spinnerInterval[kind]
}

// CycleKind advances to the next preset.
func (s *Spinner) CycleKind() {
	all := SpinnerKinds
	i := 0
	for j, k := range all {
		if k == s.kind {
			i = j
			break
		}
	}
	s.SetKind(all[(i+1)%len(all)])
}

// Update advances the spinner if its interval has elapsed.
// It is safe to call on every animation tick.
func (s *Spinner) Update(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	if now.Sub(s.last) >= s.interval {
		s.index = (s.index + 1) % len(s.frames)
		s.last = now
	}
}

// Frame returns the current spinner frame string.
func (s *Spinner) Frame() string {
	if len(s.frames) == 0 {
		return ""
	}
	return s.frames[s.index%len(s.frames)]
}
