package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type icon struct {
	x, y          int
	width, height int
	scale         float64
	zindex        int
	image         *ebiten.Image
}

func newIcon(x, y int, img *ebiten.Image) *icon {
	iconWidth, iconHeight := 100, 100
	sizeWidth, sizeHeight := img.Bounds().Dx(), img.Bounds().Dy()
	var scale float64
	if sizeWidth > sizeHeight {
		scale = float64(iconWidth) / float64(img.Bounds().Dx())
		iconHeight = int(float64(sizeHeight) * scale)
	} else {
		scale = float64(iconHeight) / float64(img.Bounds().Dy())
		iconWidth = int(float64(sizeWidth) * scale)
	}

	return &icon{
		x:      x,
		y:      y,
		width:  iconWidth,
		height: iconHeight,
		scale:  scale,
		zindex: 15,
		image:  img,
	}
}

func (i *icon) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(i.scale, i.scale)
	opts.GeoM.Translate(float64(i.x)-float64(i.width)/2, float64(i.y)-float64(i.height)/2)
	screen.DrawImage(i.image, opts)
}

func (i *icon) ZIndex() int {
	return i.zindex
}

// TODO: これらは icon.go にいるのはふさわしくない
func newHouseIcon(x, y int) *icon {
	img, _, err := image.Decode(bytes.NewReader(houseImageData))
	if err != nil {
		log.Fatal(err)
	}

	return newIcon(x, y, ebiten.NewImageFromImage(img))
}

func newBarricadeIcon(x, y int) *icon {
	img, _, err := image.Decode(bytes.NewReader(barricadeImageData))
	if err != nil {
		log.Fatal(err)
	}

	return newIcon(x, y, ebiten.NewImageFromImage(img))
}

func newBugIcon(x, y int, bugColor bugColor) *icon {
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
	bugImage := bugsImage.SubImage(rect).(*ebiten.Image)

	return newIcon(x, y, bugImage)
}
