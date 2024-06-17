package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
	_ "image/png"
)

//go:embed assets/bugs.png
var bugsImageData []byte

type bug struct {
	game          *Game
	x, y          int
	width, height int
	zindex        int
	image         *ebiten.Image

	// 画像の拡大率。
	// TODO: 本当は画像のサイズそのものを変更したほうが見た目も処理効率も良くなる。余裕があれば後々やろう。
	scale float64
}

type bugColor int

const (
	bugsRed bugColor = iota
	bugsBlue
	bugsGreen
)

func newBug(game *Game, bugColor bugColor, x, y int) *bug {
	img, _, err := image.Decode(bytes.NewReader(bugsImageData))
	if err != nil {
		log.Fatal(err)
	}

	bugsImage := ebiten.NewImageFromImage(img)
	rect := func() image.Rectangle {
		switch bugColor {
		case bugsRed:
			return redBug()
		case bugsBlue:
			return blueBug()
		case bugsGreen:
			return greenBug()
		}
		log.Fatal("invalid bug color")
		return image.Rectangle{}
	}()
	redBugImage := bugsImage.SubImage(rect).(*ebiten.Image)

	return &bug{
		game:   game,
		x:      x,
		y:      y,
		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		scale:  1,
		image:  redBugImage,
	}
}

func redBug() image.Rectangle {
	return image.Rect(1, 5, 29, 45)
}

func blueBug() image.Rectangle {
	return image.Rect(36, 4, 65, 45)
}

func greenBug() image.Rectangle {
	return image.Rect(35, 50, 66, 96)
}

// 画面中央に配置
func (h *bug) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(h.scale, h.scale)
	opts.GeoM.Translate(float64(h.x), float64(h.y))
	screen.DrawImage(h.image, opts)
}

func (h *bug) ZIndex() int {
	return h.zindex
}
