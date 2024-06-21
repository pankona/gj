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

	infoPanel *infoPanel

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

func (g *Game) RemoveBuilding(b Building) {
	for i, building := range g.buildings {
		if building == b {
			g.buildings = append(g.buildings[:i], g.buildings[i+1:]...)
			return
		}
	}
}

func (g *Game) Update() error {
	// getClickPosition の戻り値を clickHandler.HandleClick に渡す
	// これをやると登録された Clickable の OnClick が呼ばれる
	if x, y, clicked := getClickedPosition(); clicked {
		g.clickHandler.HandleClick(x, y)
	}

	// Updater を実行
	g.updateHandler.HandleUpdate()

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
	vector.DrawFilledRect(screen, screenWidth/2, eScreenHeight/2, 1, 1, color.RGBA{255, 255, 255, 255}, true)

	g.drawHandler.HandleDraw(screen)
}

const (
	screenWidth  = 1280
	screenHeight = 960
)

const (
	// infoPanel の高さを計算
	// infoPanel の高さの分だけ、ゲーム画面の中央座標が上にずれる
	// 中央座標計算のためにあらかじめここで計算しておく
	infoPanelHeight = screenHeight / 7
	eScreenHeight   = screenHeight - infoPanelHeight - 10
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
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
		newBug(g, bugsRed, screenWidth/2-50, eScreenHeight-100),
		newBug(g, bugsRed, screenWidth/2-30, eScreenHeight-100),
		newBug(g, bugsRed, screenWidth/2-10, eScreenHeight-100),
		newBug(g, bugsRed, screenWidth/2+10, eScreenHeight-100),
		newBug(g, bugsRed, screenWidth/2+30, eScreenHeight-100),
		newBug(g, bugsRed, screenWidth/2+50, eScreenHeight-100),
	}

	for _, redBug := range redBugs {
		g.drawHandler.Add(redBug)
		g.updateHandler.Add(redBug)
		g.clickHandler.Add(redBug)
	}

	//g.drawHandler.Add(newBug(g, bugsBlue, screenWidth/2, screenHeight-100))
	//g.drawHandler.Add(newBug(g, bugsGreen, screenWidth/2+50, screenHeight-100))

	// バリケードを家のすぐ下に配置
	barricadeOnDestroyFn := func(b *barricade) {
		g.drawHandler.Remove(b)
		g.clickHandler.Remove(b)
		g.RemoveBuilding(b)
	}
	barricades := []*barricade{
		newBarricade(g, screenWidth/2-105, eScreenHeight/2+80, barricadeOnDestroyFn),
		newBarricade(g, screenWidth/2, eScreenHeight/2+80, barricadeOnDestroyFn),
		newBarricade(g, screenWidth/2+105, eScreenHeight/2+80, barricadeOnDestroyFn),
	}
	for _, barricade := range barricades {
		g.drawHandler.Add(barricade)
		g.AddBuilding(barricade)
		g.clickHandler.Add(barricade)
	}

	g.AddBuilding(house)

	g.infoPanel = newInfoPanel(g, screenWidth-20, infoPanelHeight)
	g.drawHandler.Add(g.infoPanel)

	g.clickHandler.Add(house)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
