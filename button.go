package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	game *Game

	x, y          int
	width, height int
	zindex        int

	onClick func(x, y int) bool
	onDraw  func(screen *ebiten.Image, x, y, width, height int)
}

func newButton(g *Game, x, y, width, height, zindex int, clickFn func(x, y int) bool, drawFn func(screen *ebiten.Image, x, y, width, height int)) *Button {
	return &Button{
		game: g,

		x:      x,
		y:      y,
		width:  width,
		height: height,
		zindex: zindex,

		onClick: clickFn,
		onDraw:  drawFn,
	}
}

func (b *Button) OnClick(x, y int) bool {
	return b.onClick(x, y)
}

func (b *Button) IsClicked(x, y int) bool {
	return b.x <= x && x <= b.x+b.width && b.y <= y && y <= b.y+b.height
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.onDraw(screen, b.x, b.y, b.width, b.height)
}

func (b *Button) ZIndex() int {
	return b.zindex
}

func newReadyButton(g *Game) *Button {
	width, height := 100, 40
	x := screenWidth - width - 12
	y := eScreenHeight - height - 20

	return newButton(g, x, y, width, height, 1,
		func(x, y int) bool { return false },
		func(screen *ebiten.Image, x, y, width, height int) {
			const buttonMargin = 2 // 枠の幅

			// ボタンの枠を描く（白）
			vector.DrawFilledRect(screen,
				float32(x-buttonMargin), float32(y-buttonMargin),
				float32(width+2*buttonMargin), float32(height+2*buttonMargin),
				color.White, true)

			// ボタンの背景を描く（黒）
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), float32(height), color.Black, true)

			// ボタンのテキストを描く（白）
			ebitenutil.DebugPrintAt(screen, "READY", x+width/2-15, y+height/2-7)
		})
}
