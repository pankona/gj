package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

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

	// health が 0 になったときに呼ばれる関数
	onDestroy func(b *barricade)
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

// 画面中央に配置
func (b *barricade) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(b.scale, b.scale)
	opts.GeoM.Translate(float64(b.x)-float64(b.width)*b.scale/2, float64(b.y)-float64(b.height)*b.scale/2)
	screen.DrawImage(b.image, opts)
}

func (b *barricade) ZIndex() int {
	return b.zindex
}

func (b *barricade) Position() (int, int) {
	return b.x, b.y
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
		b.health = 0
		b.onDestroy(b)
	}
}

// barricade implements Clickable interface
func (b *barricade) OnClick(x, y int) {
	b.game.clickedObject = "barricade"
	// infoPanel に情報を表示する
	icon := newBarricadeIcon(80, eScreenHeight+70)
	b.game.infoPanel.setIcon(icon)
	b.game.infoPanel.setUnit(b)
}

func (b *barricade) Health() int {
	return b.health
}

func (b *barricade) IsClicked(x, y int) bool {
	w, h := b.Size()
	return b.x-w/2 <= x && x <= b.x+w/2 && b.y-h/2 <= y && y <= b.y+h/2
}
