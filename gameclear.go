package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type gameclear struct {
	game         *Game
	erapsedFrame int // このフレームを経過しないとクリックできないようにする
}

func newGameClear(g *Game) *gameclear {
	return &gameclear{
		game: g,
		// 3秒間はクリックを受け付けない
		erapsedFrame: 60 * 3,
	}
}

func (g *gameclear) Update() {
	if g.erapsedFrame > 0 {
		g.erapsedFrame--
	}
}

func (g *gameclear) OnClick(x, y int) bool {
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

func (g *gameclear) IsClicked(x, y int) bool {
	return true
}

func (g *gameclear) ZIndex() int {
	return 200
}

func (g *gameclear) Draw(screen *ebiten.Image) {
	// ゲームクリア画面の描画

	// 画面全体を半透明の黒で覆う
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{0, 0, 0, 0x90}, true)

	// とりあえず画面中央に Game Clear って出す
	ebitenutil.DebugPrintAt(screen, "Congratulation! Game Clear!", screenWidth/2-80, eScreenHeight/2-20)
	// Click to Restart って出す
	ebitenutil.DebugPrintAt(screen, "Click to Restart", screenWidth/2-50, eScreenHeight/2+10)
}
