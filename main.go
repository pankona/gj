package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	clickedPositionX, clickedPositionY int

	clickHandler *OnClickHandler
	drawHandler  *DrawHandler

	circle *Circle

	clickedObject string
}

// Game implements clickable interface
func (g *Game) OnClick() {
	g.clickedObject = ""
}

func (g *Game) IsClicked(x, y int) bool {
	return true
}

func (g *Game) ZIndex() int {
	return 0
}

func (g *Game) Update() error {
	// getClickPosition の戻り値を clickHandler.HandleClick に渡す
	// これをやると登録された Clickable の OnClick が呼ばれる
	if x, y, clicked := getClickedPosition(); clicked {
		g.clickHandler.HandleClick(x, y)
	}

	x, y, clicked := getClickedPosition()
	if clicked {
		g.clickedPositionX = x
		g.clickedPositionY = y
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// クリックされた位置を表示
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Clicked Position: (%d, %d)", g.clickedPositionX, g.clickedPositionY), 0, 0)
	// クリックされたオブジェクトを表示
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Clicked Object: %s", g.clickedObject), 0, 20)

	g.drawHandler.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 960
}

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Hello, World!")

	g := &Game{}

	// ClickHandler setup
	clickHandler := &OnClickHandler{}

	circle := &Circle{
		game: g,

		// 画面中央に配置
		x: 1280 / 2,
		y: 960 / 2,

		radius: 50,
		zindex: 1,

		image: ebiten.NewImage(1, 1),
	}

	g.clickHandler = clickHandler
	g.circle = circle

	clickHandler.Add(g)
	clickHandler.Add(circle)

	g.drawHandler = &DrawHandler{}
	g.drawHandler.Add(circle)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
