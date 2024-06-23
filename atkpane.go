package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
)

// 攻撃のためのクリックを受け止める透明なパネル
// 情報パネルを除く画面全体に敷かれ、それなりに高い ZIndex を持つ
type attackPane struct {
	game *Game

	x, y          int
	width, height int
	zindex        int
}

var (
	//go:embed assets/hand_small.png
	handSmallImageData []byte
	//go:embed assets/hand_big.png
	handBigImageData []byte
)

type smallHand struct {
	game *Game

	x, y          int
	width, height int
	zindex        int

	// 表示している時間
	// これが 0 になったら消える
	// この値は Update で減らす
	displayTime int

	cooldown    int
	erapsedTime int // 攻撃実行からの経過時間

	attackPower int

	image *ebiten.Image
}

// 同時に表示する small hand はひとつである。
// 何度も生成する必要がないので、一つだけ生成して使いまわす
var smallHandPool *smallHand

func newSmallHand(game *Game) *smallHand {
	if smallHandPool != nil {
		return smallHandPool
	}

	img, _, err := image.Decode(bytes.NewReader(handSmallImageData))
	if err != nil {
		log.Fatal(err)
	}

	h := &smallHand{
		game: game,

		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		zindex: 100,
		image:  ebiten.NewImageFromImage(img),

		cooldown: 15, // ここを短くすると連打できるようになっていく

		attackPower: 1,
	}

	smallHandPool = h

	return smallHandPool
}

func (h *smallHand) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// width と height を考慮する
	op.GeoM.Translate(float64(h.x)-float64(h.width)/2, float64(h.y)-float64(h.height)/2)
	screen.DrawImage(h.image, op)
}

func (h *smallHand) ZIndex() int {
	return h.zindex
}

// smallHand implements updater interface
func (h *smallHand) Update() {
	h.erapsedTime++
	h.displayTime--
	if h.displayTime <= 0 {
		h.game.updateHandler.Remove(h)
		h.game.drawHandler.Remove(h)
	}

	// クリックから 5 フレーム後に攻撃を実行する
	if h.erapsedTime == 5 {
		// 攻撃範囲内にいる敵に対してダメージを与える
		// ループの中で複数の h.game.enemies が減る可能性があるので、逆順でループする
		for i := len(h.game.enemies) - 1; i >= 0; i-- {
			e := h.game.enemies[i]

			ex, ey := e.Position()
			ew, eh := e.Size()

			if intersects(
				rect{h.x - h.width/2, h.y - h.height/2, h.width, h.height},
				rect{ex - ew/2, ey - eh/2, ew, eh},
			) {
				var b Damager = e.(*bug)
				b.Damage(h.attackPower)
			}
		}
	}

}

func (h *smallHand) setPosition(x, y int) {
	h.x = x
	h.y = y
}

func newAttackPane(game *Game) *attackPane {
	return &attackPane{
		game: game,

		x:      0,
		y:      0,
		width:  screenWidth,
		height: eScreenHeight,
		zindex: 100,
	}
}

// attackPane implement Clickable interface
// attackPane はクリックが下のオブジェクトに貫通する。攻撃中でも建物や敵の情報を見ることができるようにするため
func (a *attackPane) OnClick(x, y int) bool {
	hand := newSmallHand(a.game)

	// cooldown があけていなかったら攻撃を発動しない
	if hand.erapsedTime != 0 && hand.erapsedTime < hand.cooldown {
		return true
	}

	// smallHand をクリック位置に表示する
	hand.setPosition(x, y)

	// 見た目とクールダウンを一致させているが、かならずしもそうではないかも
	hand.displayTime = hand.cooldown
	hand.erapsedTime = 0

	a.game.updateHandler.Add(hand)
	a.game.drawHandler.Add(hand)

	return true
}

func (a *attackPane) IsClicked(x, y int) bool {
	return a.x <= x && x <= a.x+a.width && a.y <= y && y <= a.y+a.height
}

func (a *attackPane) ZIndex() int {
	return a.zindex
}

func (a *attackPane) RemoveAll() {
	// 初期化順序の関係で a が nil になることがある
	// TODO: みっともないので直せたら直す
	if a == nil {
		return
	}
	a.game.clickHandler.Remove(a)
}
