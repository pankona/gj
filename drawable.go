package main

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawable interface {
	Draw(screen *ebiten.Image)
	ZIndex() int
}

type DrawHandler struct {
	// List of drawable objects
	drawable []Drawable
}

func (o *DrawHandler) Add(obj Drawable) {
	o.drawable = append(o.drawable, obj)

	// Sort by ZIndex
	// ZIndex が大きいものほどあとに描画されるようにする
	// つまり ZIndex が大きいものほど後ろにくるようにソートする
	sort.Slice(o.drawable, func(i, j int) bool {
		return o.drawable[i].ZIndex() < o.drawable[j].ZIndex()
	})
}

func (o *DrawHandler) Remove(obj Drawable) {
	for i, v := range o.drawable {
		if v == obj {
			o.drawable = append(o.drawable[:i], o.drawable[i+1:]...)
			return
		}
	}
}

func (o *DrawHandler) HandleDraw(screen *ebiten.Image) {
	for _, obj := range o.drawable {
		obj.Draw(screen)
	}
}

func (o *DrawHandler) Clear() {
	o.drawable = []Drawable{}
}
