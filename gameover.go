package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type gameover struct {
	game *Game
}

func newGameover(g *Game) *gameover {
	return &gameover{
		game: g,
	}
}

func (g *gameover) OnClick(x, y int) bool {
	// ゲームリセットする

	g.game.Reset()

	// ゲームオーバー画面を削除
	g.game.clickHandler.Remove(g)
	g.game.drawHandler.Remove(g)

	return false
}

func (g *gameover) IsClicked(x, y int) bool {
	return true
}

func (g *gameover) ZIndex() int {
	return 200
}

// gameover implements drawable
func (g *gameover) Draw(screen *ebiten.Image) {
	// ゲームオーバー画面の描画

	// 画面全体を半透明の黒で覆う
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{0, 0, 0, 0x90}, true)

	// とりあえず画面中央に GameOver って出す
	ebitenutil.DebugPrintAt(screen, "Game Over!", screenWidth/2-30, eScreenHeight/2-20)
	// Click to Restart って出す
	ebitenutil.DebugPrintAt(screen, "Click to Restart", screenWidth/2-50, eScreenHeight/2+10)
}
