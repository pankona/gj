package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	game   *Game
	zindex int
	onDraw func(screen *ebiten.Image)
}

func newButton(g *Game, zindex int, drawFn func(screen *ebiten.Image)) *Button {
	return &Button{
		game:   g,
		zindex: zindex,
		onDraw: drawFn,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.onDraw(screen)
}

func (b *Button) ZIndex() int {
	return b.zindex
}

func newReadyButton(g *Game) *Button {
	return newButton(g, 1, func(screen *ebiten.Image) {
		buttonWidth, buttonHeight := 100, 40
		buttonX := screenWidth - buttonWidth - 12
		buttonY := eScreenHeight - buttonHeight - 20

		const buttonMargin = 2 // 枠の幅

		// ボタンの枠を描く（白）
		vector.DrawFilledRect(screen,
			float32(buttonX-buttonMargin), float32(buttonY-buttonMargin),
			float32(buttonWidth+2*buttonMargin), float32(buttonHeight+2*buttonMargin),
			color.White, true)

		// ボタンの背景を描く（黒）
		vector.DrawFilledRect(screen, float32(buttonX), float32(buttonY), float32(buttonWidth), float32(buttonHeight), color.Black, true)

		// ボタンのテキストを描く（白）
		ebitenutil.DebugPrintAt(screen, "READY", buttonX+buttonWidth/2-15, buttonY+buttonHeight/2-7)
	})

}
