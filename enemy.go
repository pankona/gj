package main

type Enemy interface {
	Position() (int, int)
	Size() (int, int)
	Name() string

	Drawable
	Clickable
	Updater
}

func (g *Game) AddEnemy(e Enemy) {
	g.enemies = append(g.enemies, e)
}

func (g *Game) RemoveEnemy(e Enemy) {
	for i, enemy := range g.enemies {
		if enemy == e {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			return
		}
	}
}
