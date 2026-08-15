package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/richavery/donk-cli/internal/shader"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type visualizer struct {
	shaders    []string
	shaderIdx  int
	cursorPos  image.Point
	prevPos    image.Point
	cursorTimer float64
}

func (v *visualizer) Update() error {
	v.cursorTimer += 1.0 / 60.0
	if v.cursorTimer > 2.0 {
		v.cursorTimer = 0
		v.prevPos = v.cursorPos
		v.cursorPos = image.Point{
			X: 100 + int(600*((v.cursorTimer*0.5))),
			Y: 100 + int(400*((v.cursorTimer*0.3))),
		}
	}
	return nil
}

func (v *visualizer) Draw(screen *ebiten.Image) {
	// Dark background
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{0x1a, 0x1a, 0x2e, 0xff})
	
	// Draw cursor
	ebitenutil.DrawRect(screen, float64(v.cursorPos.X), float64(v.cursorPos.Y), 20, 20, color.White)
	
	// Draw shader info
	msg := fmt.Sprintf("DONK Shader Visualizer\nShader: %s\nPress SPACE to cycle shaders", v.shaders[v.shaderIdx])
	ebitenutil.DebugPrint(screen, msg)
}

func (v *visualizer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func listShaders() []string {
	entries, err := os.ReadDir("internal/shader")
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.Type().IsRegular() && len(e.Name()) > 5 && e.Name()[:len(e.Name())-5] == ".glsl" {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return []string{"cursor_warp.glsl", "cursor_sweep.glsl", "cursor_tail.glsl"}
	}
	return names
}

func main() {
	shaders := listShaders()
	if len(shaders) == 0 {
		log.Fatal("no shaders found")
	}
	
	// Try to load default shader
	_, err := shader.Read("cursor_warp.glsl")
	if err != nil {
		log.Printf("warning: failed to read embedded shader: %v", err)
	}
	
	v := &visualizer{
		shaders:   shaders,
		shaderIdx: 0,
		cursorPos: image.Point{400, 300},
		prevPos:   image.Point{400, 300},
	}
	
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("DONK Shader Visualizer")
	
	if err := ebiten.RunGame(v); err != nil {
		log.Fatal(err)
	}
}
