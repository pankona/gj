package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// パネルの枠を表示するための構造体
type infoPanel struct {
	game *Game

	x, y          int
	width, height int
	zindex        int

	icon *icon
	unit infoer
}

func newInfoPanel(g *Game, w, h int) *infoPanel {
	bottomMargin := 10
	return &infoPanel{
		game:   g,
		x:      screenWidth/2 - w/2,
		y:      screenHeight - h - bottomMargin,
		width:  w,
		height: h,
		zindex: 10,
	}
}

func (p *infoPanel) setIcon(i *icon) {
	p.game.drawHandler.Remove(p.icon)

	p.icon = i
	if p.icon == nil {
		return
	}
}

type infoer interface {
	Name() string
	Health() int
}

func (p *infoPanel) setUnit(u infoer) {
	p.unit = u
}

func (p *infoPanel) Remove(u infoer) {
	if p.unit == u {
		p.unit = nil
	}
}

func (p *infoPanel) Draw(screen *ebiten.Image) {
	// 枠を描画
	strokeWidth := float32(2)
	vector.StrokeLine(screen, float32(p.x), float32(p.y), float32(p.x+p.width), float32(p.y), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(p.x), float32(p.y), float32(p.x), float32(p.y+p.height), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(p.x+p.width), float32(p.y), float32(p.x+p.width), float32(p.y+p.height), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(p.x), float32(p.y+p.height), float32(p.x+p.width), float32(p.y+p.height), strokeWidth, color.White, true)

	// TODO: アイコンを描画

	// ユニット名とHPを描画
	if p.unit == nil {
		return
	}
	p.icon.Draw(screen)
	name, health := p.unit.Name(), p.unit.Health()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", name), p.x+100+40, p.y+30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d", health), p.x+100+40, p.y+50)
}

func (p *infoPanel) ZIndex() int {
	return p.zindex
}
