package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
	_ "image/png"
)

//go:embed assets/tower.png
var towerImageData []byte

type tower struct {
	game *Game

	x, y          int
	width, height int
	zindex        int
	image         *ebiten.Image

	health int

	// 画像の拡大率。
	// 1以外を指定する場合は元画像のサイズをそもそも変更できないか検討すること
	scale float64

	// health が 0 になったときに呼ばれる関数
	onDestroy func(b *tower)

	// この建物が他の建物と重なっているかどうか (建築確定前に用いるフラグ)
	isOverlapping bool
}

func newTower(game *Game, x, y int, onDestroy func(b *tower)) *tower {
	img, _, err := image.Decode(bytes.NewReader(towerImageData))
	if err != nil {
		log.Fatal(err)
	}

	h := &tower{
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

// 画面中央に配置
func (b *tower) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(b.scale, b.scale)
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

func (b *tower) ZIndex() int {
	return b.zindex
}

func (b *tower) Position() (int, int) {
	return b.x, b.y
}

func (b *tower) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

func (b *tower) Size() (int, int) {
	return int(float64(b.width) * b.scale), int(float64(b.height) * b.scale)
}

func (b *tower) Name() string {
	return "Tower"
}

func (b *tower) Damage(d int) {
	if b.health <= 0 {
		return
	}

	b.health -= d
	if b.health <= 0 {
		b.health = 0
		b.onDestroy(b)
	}
}

// tower implements Clickable interface
func (b *tower) OnClick(x, y int) bool {
	b.game.clickedObject = "tower"

	// infoPanel に情報を表示する

	// TODO: ClearButtons は呼び出し側でやるんじゃなくて infoPanel 側のどっかでやるべきかな
	b.game.infoPanel.ClearButtons()
	icon := newTowerIcon(80, eScreenHeight+70)
	b.game.infoPanel.setIcon(icon)
	b.game.infoPanel.setUnit(b)

	return false
}

func newTowerIcon(x, y int) *icon {
	img, _, err := image.Decode(bytes.NewReader(towerImageData))
	if err != nil {
		log.Fatal(err)
	}

	return newIcon(x, y, ebiten.NewImageFromImage(img))
}

func (b *tower) Health() int {
	return b.health
}

func (b *tower) IsClicked(x, y int) bool {
	w, h := b.Size()
	return b.x-w/2 <= x && x <= b.x+w/2 && b.y-h/2 <= y && y <= b.y+h/2
}

func (b *tower) SetOverlap(overlap bool) {
	b.isOverlapping = overlap
}

func (b *tower) IsOverlap() bool {
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

func (b *tower) Cost() int {
	return CostTowerBuild
}
