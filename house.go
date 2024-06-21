package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
	_ "image/png"
)

//go:embed assets/house.png
var houseImageData []byte

type house struct {
	game *Game

	x, y          int // 画面中央に配置するので初期化時に値をもらう必要はない
	width, height int // 画像サイズをそのまま使うので初期化時に値をもらう必要はない
	zindex        int // これも適当に調整するので初期化時に値をもらう必要はない
	image         *ebiten.Image

	health int

	// 画像の拡大率。
	// TODO: 本当は画像のサイズそのものを変更したほうが見た目も処理効率も良くなる。余裕があれば後々やろう。
	scale float64
}

func newHouse(game *Game) *house {
	img, _, err := image.Decode(bytes.NewReader(houseImageData))
	if err != nil {
		log.Fatal(err)
	}

	h := &house{
		game: game,

		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		scale:  0.5,

		health: 100,

		image: ebiten.NewImageFromImage(img),
	}

	h.x = screenWidth / 2
	h.y = eScreenHeight / 2

	return h
}

// 画面中央に配置
func (h *house) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(h.scale, h.scale)
	opts.GeoM.Translate(float64(h.x)-float64(h.width)*h.scale/2, float64(h.y)-float64(h.height)*h.scale/2)
	screen.DrawImage(h.image, opts)
}

func (h *house) ZIndex() int {
	return h.zindex
}

func (h *house) Position() (int, int) {
	// 中央の座標を返す
	return h.x, h.y
}

func (h *house) Size() (int, int) {
	return int(float64(h.width) * h.scale), int(float64(h.height) * h.scale)
}

func (h *house) Name() string {
	return "House"
}

func (h *house) Damage(d int) {
	if h.health <= 0 {
		return
	}

	h.health -= d
	if h.health <= 0 {
		h.health = 0
	}
}

// house implements Clickable interface
func (h *house) OnClick() {
	h.game.clickedObject = "House"
	// infoPanel に情報を表示する
	icon := newHouseIcon(80, eScreenHeight+70)
	h.game.infoPanel.setIcon(icon)
	h.game.infoPanel.setUnit(h)
}

func (h *house) Health() int {
	return h.health
}

func (h *house) IsClicked(x, y int) bool {
	width, height := h.Size()
	return h.x-width/2 <= x && x <= h.x+width/2 && h.y-height/2 <= y && y <= h.y+height/2
}
