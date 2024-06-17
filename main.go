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

	clickedObject string
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

	g := &Game{
		clickHandler: &OnClickHandler{},
		drawHandler:  &DrawHandler{},
	}

	// 最初のシーンをセットアップする
	// とりあえずいきなりゲームが始まるとする。
	// TODO: まずタイトルバックを表示して、その後にゲーム画面に遷移するようにする
	g.drawHandler.Add(newHouse(g))
	g.drawHandler.Add(newReadyButton(g))

	//とりあえずいったん虫を画面の下部に配置
	screenWidth, screenHeight := 1280, 960
	g.drawHandler.Add(newBug(g, bugsRed, screenWidth/2-50, screenHeight-100))
	g.drawHandler.Add(newBug(g, bugsBlue, screenWidth/2, screenHeight-100))
	g.drawHandler.Add(newBug(g, bugsGreen, screenWidth/2+50, screenHeight-100))

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
