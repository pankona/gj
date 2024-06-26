package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "embed"
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

	readyButton *Button
}

func newBuildPane(game *Game) *buildPane {
	okButton := newButton(game, screenWidth-200-22, eScreenHeight-80, 100, 50, 110,
		func(x, y int) bool {
			if game.buildCandidate == nil {
				return true
			}

			bcX, bcY := game.buildCandidate.Position()
			if bcX == 0 && bcY == 0 {
				// まだ場所が決まっていない場合はボタンを無効にする
				return true
			}

			// 建築不可能な場所を指定していた場合は何もしない
			if game.buildCandidate.IsOverlap() {
				return false
			}

			// クリックされたら建築を確定する
			game.AddBuilding(game.buildCandidate)
			game.clickHandler.Add(game.buildCandidate)
			game.updateHandler.Add(game.buildCandidate)

			// クレジットを減らす
			game.credit -= game.buildCandidate.Cost()

			// buildCandidate は次の建築のために初期化する
			game.buildCandidate = nil
			game.infoPanel.drawDescriptionFn = nil

			getAudioPlayer().play(soundDon)

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
			// 置けない場所に建築しようとした場合はボタンをグレーアウトする
			if game.buildCandidate.IsOverlap() {
				drawGrayRect(screen, x, y, width, height)
			} else {
				drawRect(screen, x, y, width, height)
			}
			ebitenutil.DebugPrintAt(screen, "BUILD IT!", x+width/2-25, y+height/2-8)
		})

	cancelButton := newButton(game, screenWidth-100-12, eScreenHeight-80, 100, 50, 110,
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
			game.infoPanel.drawDescriptionFn = nil

			getAudioPlayer().play(soundChoice)

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

	// 他の建築物と重なっているかどうか判定してフラグをセットする
	a.game.buildCandidate.SetOverlap(a.game.buildCandidate.IsOverlap())

	return true
}

func (a *buildPane) IsClicked(x, y int) bool {
	return a.x <= x && x <= a.x+a.width && a.y <= y && y <= a.y+a.height
}

func (a *buildPane) ZIndex() int {
	return a.zindex
}

// ウェーブフェーズに繊維するとき、建築フェーズのパネルを取り去る
// このとき、ClickHandler や DrawHandler に入れたボタンも取り去る
func (p *buildPane) RemoveAll() {
	p.game.clickHandler.Remove(p.okButton)
	p.game.clickHandler.Remove(p.cancelButton)
	p.game.clickHandler.Remove(p.readyButton)
	p.game.drawHandler.Remove(p)
	p.game.clickHandler.Remove(p)
}
