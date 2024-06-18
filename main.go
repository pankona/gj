package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	clickedPositionX, clickedPositionY int

	clickHandler  *OnClickHandler
	drawHandler   *DrawHandler
	updateHandler *UpdateHandler

	// 建物のリスト
	buildings []Building

	clickedObject string
}

type Building interface {
	Position() (int, int)
	Size() (int, int)
	Name() string
}

func (g *Game) AddBuilding(b Building) {
	g.buildings = append(g.buildings, b)
}

func (g *Game) Update() error {
	// getClickPosition の戻り値を clickHandler.HandleClick に渡す
	// これをやると登録された Clickable の OnClick が呼ばれる
	if x, y, clicked := getClickedPosition(); clicked {
		g.clickHandler.HandleClick(x, y)
	}

	// Updater を実行
	g.updateHandler.Update()

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

	// 画面中央に点を表示 (debug)
	vector.DrawFilledRect(screen, screenWidth/2, screenHeight/2, 1, 1, color.RGBA{255, 255, 255, 255}, true)

	g.drawHandler.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 960
}

const (
	screenWidth  = 1280
	screenHeight = 960
)

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Hello, World!")

	g := &Game{
		clickHandler:  &OnClickHandler{},
		drawHandler:   &DrawHandler{},
		updateHandler: &UpdateHandler{},
	}

	// 最初のシーンをセットアップする
	// とりあえずいきなりゲームが始まるとする。
	// TODO: まずタイトルバックを表示して、その後にゲーム画面に遷移するようにする
	house := newHouse(g)
	g.drawHandler.Add(house)
	g.drawHandler.Add(newReadyButton(g))

	//とりあえずいったん虫を画面の下部に配置
	redBugs := []*bug{
		newBug(g, bugsRed, screenWidth/2-50, screenHeight-100),
		newBug(g, bugsRed, screenWidth/2-30, screenHeight-100),
		newBug(g, bugsRed, screenWidth/2-10, screenHeight-100),
		newBug(g, bugsRed, screenWidth/2+10, screenHeight-100),
		newBug(g, bugsRed, screenWidth/2+30, screenHeight-100),
		newBug(g, bugsRed, screenWidth/2+50, screenHeight-100),
	}

	for _, redBug := range redBugs {
		g.drawHandler.Add(redBug)
		g.updateHandler.Add(redBug)
	}

	//g.drawHandler.Add(newBug(g, bugsBlue, screenWidth/2, screenHeight-100))
	//g.drawHandler.Add(newBug(g, bugsGreen, screenWidth/2+50, screenHeight-100))

	// 建物一覧に登録
	g.AddBuilding(house)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
