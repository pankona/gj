package main

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "embed"
	_ "image/png"
)

//go:embed assets/barricade.png
var barricadeImageData []byte

type barricade struct {
	game *Game

	x, y          int
	width, height int
	zindex        int
	image         *ebiten.Image

	health int

	// 画像の拡大率。
	// 1以外を指定する場合は元画像のサイズをそもそも変更できないか検討すること
	scale float64

	// 壊れたときのアニメーションを制御するための変数
	deadAnimationDuration int

	// health が 0 になったときに呼ばれる関数
	onDestroy func(b *barricade)

	// この建物が他の建物と重なっているかどうか (建築確定前に用いるフラグ)
	isOverlapping bool
}

func newBarricade(game *Game, x, y int, onDestroy func(b *barricade)) *barricade {
	img, _, err := image.Decode(bytes.NewReader(barricadeImageData))
	if err != nil {
		log.Fatal(err)
	}

	h := &barricade{
		game: game,

		x:      x,
		y:      y,
		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		scale:  1,

		health: 100,

		image: ebiten.NewImageFromImage(img),

		onDestroy: onDestroy,
	}

	return h
}

func (b *barricade) Update() {
	if b.health <= 0 {
		b.deadAnimationDuration++
		if b.deadAnimationDuration >= deadAnimationTotalFrame {
			b.onDestroy(b)
		}
		return
	}
}

// 画面中央に配置
func (b *barricade) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	if b.health <= 0 {
		// 死亡時のアニメーションを行う
		// ぺちゃんこになるように縮小する
		scale := 1.0 - float64(b.deadAnimationDuration)/deadAnimationTotalFrame
		if scale < 0 {
			scale = 0
		}

		opts.GeoM.Translate(0, float64(-b.height))
		opts.GeoM.Scale(1, scale)
		opts.GeoM.Translate(0, float64(b.height))
	} else {
		opts.GeoM.Scale(b.scale, b.scale)
	}
	opts.GeoM.Translate(float64(b.x)-float64(b.width)*b.scale/2, float64(b.y)-float64(b.height)*b.scale/2)

	// 他の建物と重なっている場合は赤くする
	if b.isOverlapping {
		opts.ColorScale.Scale(1, 0, 0, 1)
	} else if b.game.buildCandidate == b {
		// 建築確定前は暗い色で建物を描画する
		opts.ColorScale.Scale(0.5, 0.5, 0.5, 1)
	}

	screen.DrawImage(b.image, opts)
}

func (b *barricade) ZIndex() int {
	return b.zindex
}

func (b *barricade) Position() (int, int) {
	return b.x, b.y
}

func (b *barricade) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

func (b *barricade) Size() (int, int) {
	return int(float64(b.width) * b.scale), int(float64(b.height) * b.scale)
}

func (b *barricade) Name() string {
	return "Barricade"
}

func (b *barricade) Damage(d int) {
	if b.health <= 0 {
		return
	}

	b.health -= d
	if b.health <= 0 {
		getAudioPlayer().play(soundKuzureru)
		b.health = 0
	}
}

// barricade implements Clickable interface
func (b *barricade) OnClick(x, y int) bool {
	if b.game.buildCandidate != nil {
		// 建築予定のものを持っているときには何もしない
		return false
	}

	b.game.clickedObject = "barricade"
	getAudioPlayer().play(soundChoice)

	// infoPanel に情報を表示する

	// TODO: ClearButtons は呼び出し側でやるんじゃなくて infoPanel 側のどっかでやるべきかな
	b.game.infoPanel.ClearButtons()
	icon := newBarricadeIcon(80, eScreenHeight+70)
	b.game.infoPanel.setIcon(icon)
	b.game.infoPanel.setUnit(b)
	b.game.infoPanel.drawDescriptionFn = func(screen *ebiten.Image, x, y int) {
		ebitenutil.DebugPrintAt(screen, "I am Barricade!", x, y)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cost: $%d", b.Cost()), x, y+20)
		// 敵の進行を邪魔するという説明を記載する
		ebitenutil.DebugPrintAt(screen, "Blocks enemy's advance!", x, y+40)
	}

	return false
}

func (b *barricade) Health() int {
	return b.health
}

func (b *barricade) IsClicked(x, y int) bool {
	w, h := b.Size()
	return b.x-w/2 <= x && x <= b.x+w/2 && b.y-h/2 <= y && y <= b.y+h/2
}

func (b *barricade) SetOverlap(overlap bool) {
	b.isOverlapping = overlap
}

func (b *barricade) IsOverlap() bool {
	// 他の建物と重なっているかどうかを判定する
	for _, building := range b.game.buildings {
		if building == b {
			continue
		}

		bx, by := building.Position()
		bw, bh := building.Size()

		if intersects(
			rect{b.x - b.width/2, b.y - b.height/2, b.width, b.height},
			rect{bx - bw/2, by - bh/2, bw, bh},
		) {
			return true
		}
	}

	return false
}

func (b *barricade) Cost() int {
	return CostBarricadeBuild
}
