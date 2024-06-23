package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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
