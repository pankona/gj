package main

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// 建築のためのクリックを受け止める透明なパネル
// 情報パネルを除く画面全体に敷かれ、それなりに高い ZIndex を持つ
type buildPane struct {
	game *Game

	x, y          int
	width, height int
	zindex        int

	okButton     *Button
	cancelButton *Button
}

func newBuildPane(game *Game) *buildPane {
	okButton := newButton(game, screenWidth-200-22, eScreenHeight-130, 100, 50, 110,
		func(x, y int) bool {
			if game.buildCandidate == nil {
				return true
			}

			bcX, bcY := game.buildCandidate.Position()
			if bcX == 0 && bcY == 0 {
				// まだ場所が決まっていない場合はボタンを無効にする
				return true
			}

			// クリックされたら建築を確定する

			game.AddBuilding(game.buildCandidate)
			game.buildCandidate = nil

			return false
		},
		func(screen *ebiten.Image, x, y, width, height int) {
			if game.buildCandidate == nil {
				return
			}
			bcX, bcY := game.buildCandidate.Position()
			if bcX == 0 && bcY == 0 {
				// まだ場所が決まっていない場合はボタンを無効にする
				return
			}

			drawRect(screen, x, y, width, height)
			ebitenutil.DebugPrintAt(screen, "OK", x+width/2-10, y+height/2-8)
		})

	cancelButton := newButton(game, screenWidth-100-12, eScreenHeight-130, 100, 50, 110,
		func(x, y int) bool {
			if game.buildCandidate == nil {
				return true
			}
			bcX, bcY := game.buildCandidate.Position()
			if bcX == 0 && bcY == 0 {
				// まだ場所が決まっていない場合はボタンを無効にする
				return true
			}

			// クリックされたら建築をキャンセルする

			game.drawHandler.Remove(game.buildCandidate)
			game.buildCandidate = nil

			return false
		},
		func(screen *ebiten.Image, x, y, width, height int) {
			if game.buildCandidate == nil {
				return
			}
			bcX, bcY := game.buildCandidate.Position()
			if bcX == 0 && bcY == 0 {
				// まだ場所が決まっていない場合はボタンを無効にする
				return
			}

			drawRect(screen, x, y, width, height)
			ebitenutil.DebugPrintAt(screen, "Cancel", x+width/2-20, y+height/2-8)
		})

	game.clickHandler.Add(okButton)
	game.clickHandler.Add(cancelButton)

	return &buildPane{
		game: game,

		x:      0,
		y:      0,
		width:  screenWidth,
		height: eScreenHeight,
		zindex: 100,

		okButton:     okButton,
		cancelButton: cancelButton,
	}
}

func (a *buildPane) Draw(screen *ebiten.Image) {
	a.okButton.Draw(screen)
	a.cancelButton.Draw(screen)
}

// buildPane implement Clickable interface
// buildPane はクリックが下のオブジェクトに貫通する。建築中でも建物や敵の情報を見ることができるようにするため
func (a *buildPane) OnClick(x, y int) bool {
	if a.game.buildCandidate == nil {
		return true
	}

	if !a.game.drawHandler.Lookup(a.game.buildCandidate) {
		a.game.drawHandler.Add(a.game.buildCandidate)
	}

	a.game.buildCandidate.SetPosition(x, y)

	return true
}

func (a *buildPane) IsClicked(x, y int) bool {
	return a.x <= x && x <= a.x+a.width && a.y <= y && y <= a.y+a.height
}

func (a *buildPane) ZIndex() int {
	return a.zindex
}
