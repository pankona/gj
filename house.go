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
	game          *Game
	width, height int
	zindex        int
	image         *ebiten.Image

	// 画像の拡大率。
	// TODO: 本当は画像のサイズそのものを変更したほうが見た目も処理効率も良くなる。余裕があれば後々やろう。
	scale float64
}

func NewHouse(game *Game) *house {
	img, _, err := image.Decode(bytes.NewReader(houseImageData))
	if err != nil {
		log.Fatal(err)
	}

	return &house{
		game:   game,
		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		scale:  0.5,
		image:  ebiten.NewImageFromImage(img),
	}
}

// 画面中央に配置
func (h *house) Draw(screen *ebiten.Image) {
	// ウィンドウサイズを取得
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// 画像を中央に配置するための座標を計算
	x := (screenWidth - int(float64(h.width)*h.scale)) / 2
	y := (screenHeight - int(float64(h.height)*h.scale)) / 2

	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(h.scale, h.scale)
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(h.image, opts)
}

func (h *house) ZIndex() int {
	return h.zindex
}
