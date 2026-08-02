package anim

import "testing"

func TestGalleryTabLineBracketsActive(t *testing.T) {
	got := GalleryTabLine(1)
	want := "Harmonica Spring  [Gradient Splash]  Doom Fire  Blinking Eyes  Char Cycling"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGalleryTabLineFirst(t *testing.T) {
	got := GalleryTabLine(0)
	want := "[Harmonica Spring]  Gradient Splash  Doom Fire  Blinking Eyes  Char Cycling"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGalleryTabLineLast(t *testing.T) {
	got := GalleryTabLine(len(GalleryNames) - 1)
	want := "Harmonica Spring  Gradient Splash  Doom Fire  Blinking Eyes  [Char Cycling]"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
