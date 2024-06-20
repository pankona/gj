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

	// house の HP を表示
	for _, building := range g.buildings {
		if building.Name() == "house" {
			h := building.(*house)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("House HP: %d", h.health), 0, 40)
		}
	}

	g.drawHandler.Draw(screen)
}

const (
	screenWidth  = 1280
	screenHeight = 960
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

	// バリケードを家のすぐ下に配置
	barricades := []*barricade{
		newBarricade(g, screenWidth/2-105, screenHeight/2+80),
		newBarricade(g, screenWidth/2, screenHeight/2+80),
		newBarricade(g, screenWidth/2+105, screenHeight/2+80),
	}
	for _, barricade := range barricades {
		g.drawHandler.Add(barricade)
		g.AddBuilding(barricade)
	}

	g.AddBuilding(house)

	g.drawHandler.Add(newInfoPanel(screenWidth-20, screenHeight/7))

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// パネルの枠を表示するための構造体
type infoPanel struct {
	x, y          int
	width, height int
	zindex        int
}

func newInfoPanel(w, h int) *infoPanel {
	bottomMargin := 10
	return &infoPanel{
		x:      screenWidth/2 - w/2,
		y:      screenHeight - h - bottomMargin,
		width:  w,
		height: h,
		zindex: 10,
	}
}

func (p *infoPanel) Draw(screen *ebiten.Image) {
	// 枠を描画
	strokeWidth := float32(2)
	vector.StrokeLine(screen, float32(p.x), float32(p.y), float32(p.x+p.width), float32(p.y), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(p.x), float32(p.y), float32(p.x), float32(p.y+p.height), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(p.x+p.width), float32(p.y), float32(p.x+p.width), float32(p.y+p.height), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(p.x), float32(p.y+p.height), float32(p.x+p.width), float32(p.y+p.height), strokeWidth, color.White, true)
}

func (p *infoPanel) ZIndex() int {
	return p.zindex
}
