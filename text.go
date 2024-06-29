package main

import (
	"image/color"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

/*
こんなふうに使う

clr := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
scaleX, scaleY = float64(5), float64(5)
drawText(screen, "HOUSE DEFENSE OPERATION!", screenWidth-750, 100, scaleX, scaleY, clr)
*/
func drawText(screen *ebiten.Image, t string, x, y int, scaleX, scaleY float64, clr color.RGBA) {
	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(clr)
	textOp.GeoM.Scale(scaleX, scaleY)
	textOp.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, t, text.NewGoXFace(bitmapfont.Face), textOp)
}
