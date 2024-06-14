package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Circle struct {
	game         *Game
	x, y, radius float64
	zindex       int

	image *ebiten.Image
}

func (c *Circle) OnClick() {
	c.game.clickedObject = "circle"
}

func (c *Circle) IsClicked(x, y int) bool {
	return math.Pow(c.x-float64(x), 2)+math.Pow(c.y-float64(y), 2) <= math.Pow(c.radius, 2)
}

func (c *Circle) ZIndex() int {
	return c.zindex
}

// circle implement draw
func (c *Circle) Draw(screen *ebiten.Image) {
	path := vector.Path{}
	path.MoveTo(float32(c.x+c.radius), float32(c.y))
	for i := 1; i <= 360; i++ {
		angle := float64(i) * (3.14159265 / 180)
		x := c.x + c.radius*math.Cos(angle)
		y := c.y + c.radius*math.Sin(angle)
		path.LineTo(float32(x), float32(y))
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	c.image.Fill(color.White)
	op := &ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vs, is, c.image, op)
}
