package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type background struct {
	game   *Game
	zindex int
}

func newBackground(g *Game) *background {
	return &background{
		game:   g,
		zindex: -1,
	}
}

// やや濃い緑で塗りつぶす
func (b *background) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x45, 0x00, 0xff})
}

func (b *background) ZIndex() int {
	return b.zindex
}
