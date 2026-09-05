package audio

import (
	"log/slog"

	tea "charm.land/bubbletea/v2"
)

// NativeAudioBackend sends audio notifications using the native OS audio system.
// The actual delivery function is supplied per-platform via defaultAudioFunc;
// on unsupported platforms it is a no-op. Selection logic avoids this backend there
// and uses a terminal-based backend instead, so this is only a safety net.
// See NativeSupported.
type NativeAudioBackend struct {
	// audioFunc is the function used to send audio (swappable for testing).
	audioFunc func(title, message, audioType string) error
}

// NewNativeAudioBackend creates a new native audio backend.
func NewNativeAudioBackend() *NativeAudioBackend {
	return &NativeAudioBackend{
		audioFunc: defaultAudioFunc,
	}
}

// Play returns a command that sends audio using the native OS audio system.
func (b *NativeAudioBackend) Play(a Audio) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("Sending native audio", "title", a.Title, "message", a.Message, "type", a.Type)

		if err := b.audioFunc(a.Title, a.Message, a.Type); err != nil {
			slog.Error("Failed to send audio", "error", err)
		} else {
			slog.Debug("Audio sent successfully")
		}

		return nil
	}
}

// SetAudioFunc allows replacing the audio function for testing.
func (b *NativeAudioBackend) SetAudioFunc(fn func(title, message, audioType string) error) {
	b.audioFunc = fn
}

// ResetAudioFunc resets the audio function to the default.
func (b *NativeAudioBackend) ResetAudioFunc() {
	b.audioFunc = defaultAudioFunc
}

// defaultAudioFunc is a no-op fallback for unsupported platforms.
var defaultAudioFunc = func(title, message, audioType string) error { return nil }
