package main

type Building interface {
	Position() (int, int)
	SetPosition(int, int)
	Size() (int, int)
	Name() string

	SetOverlap(bool)
	IsOverlap() bool

	Cost() int

	Clickable
	Drawable
}

func (g *Game) AddBuilding(b Building) {
	g.buildings = append(g.buildings, b)
}

func (g *Game) RemoveBuilding(b Building) {
	for i, building := range g.buildings {
		if building == b {
			g.buildings = append(g.buildings[:i], g.buildings[i+1:]...)
			return
		}
	}
}
