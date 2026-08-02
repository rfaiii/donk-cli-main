// Package anim provides animation gallery metadata and registries
// ported from the Rust donk-anim crate.
package anim

// GalleryName is the display name for an animation or scene.
type GalleryName string

// GalleryNames is the ordered list of available gallery items.
// This is ported from Rust `donk-anim/src/gallery.rs`.
var GalleryNames = []GalleryName{
	"Harmonica Spring",
	"Gradient Splash",
	"Doom Fire",
	"Blinking Eyes",
	"Char Cycling",
}

// GalleryTabLine returns a tab-style string with the active item bracketed.
// This matches the behavior of Rust `gallery_tab_line`.
func GalleryTabLine(active int) string {
	parts := make([]string, len(GalleryNames))
	for i, name := range GalleryNames {
		if i == active {
			parts[i] = "[" + string(name) + "]"
		} else {
			parts[i] = string(name)
		}
	}
	return joinTabParts(parts)
}

func joinTabParts(parts []string) string {
	var b []byte
	for i, p := range parts {
		if i > 0 {
			b = append(b, ' ', ' ')
		}
		b = append(b, p...)
	}
	return string(b)
}
