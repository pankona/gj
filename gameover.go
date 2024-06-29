package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type gameover struct {
	game *Game

	erapsedFrame int // このフレームを経過しないとクリックできないようにする
}

func newGameover(g *Game) *gameover {
	return &gameover{
		game:         g,
		erapsedFrame: 60 * 3,
	}
}

func (g *gameover) Update() {
	if g.erapsedFrame > 0 {
		g.erapsedFrame--
	}
}

func (g *gameover) OnClick(x, y int) bool {
	// ちょっとの間クリックを受け付けないようにする
	if g.erapsedFrame > 0 {
		return false
	}

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

	// 画面中央に負けた感のあるメッセージを出す
	drawText(screen, "You lose! House destroyed...", screenWidth/2-400, eScreenHeight/2-100, 5, 5, color.RGBA{0xff, 0xff, 0xff, 0xff})
	// Click to Restart って出す
	drawText(screen, "Click to Restart", screenWidth/2-230, eScreenHeight/2+100, 5, 5, color.RGBA{0xff, 0xff, 0xff, 0xff})
}
