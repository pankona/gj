package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type instruction struct {
	game *Game

	erapsedFrame int
	zindex       int

	x, y int
	text string
}

func newInstruction(game *Game, text string, x, y int) *instruction {
	return &instruction{
		game: game,

		x:    x,
		y:    y,
		text: text,

		zindex: 100,
	}
}

func (i *instruction) Draw(screen *ebiten.Image) {
	// click me to open build menu と、家の下 (画面中央下) に表示する
	// 点滅させる
	i.erapsedFrame++
	if i.erapsedFrame%240 < 120 {
		ebitenutil.DebugPrintAt(screen, i.text, i.x, i.y)
	}
}

func (i *instruction) ZIndex() int {
	return i.zindex
}
