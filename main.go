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

	// 敵のリスト
	enemies []Enemy

	infoPanel *infoPanel

	// 建築対象としていったん保持されているオブジェクト
	buildCandidate Building

	clickedObject string
}

type Enemy interface {
	Position() (int, int)
	Size() (int, int)
	Name() string
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

type Building interface {
	Position() (int, int)
	SetPosition(int, int)
	Size() (int, int)
	Name() string

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

	//bugDestroyFn := func(b *bug) {
	//	g.drawHandler.Remove(b)
	//	g.updateHandler.Remove(b)
	//	g.clickHandler.Remove(b)
	//	g.RemoveEnemy(b)
	//	g.infoPanel.Remove(b)
	//}

	//とりあえずいったん虫を画面の下部に配置
	redBugs := []*bug{
		//newBug(g, bugsRed, screenWidth/2-50, eScreenHeight-100, bugDestroyFn),
		//newBug(g, bugsRed, screenWidth/2-30, eScreenHeight-100, bugDestroyFn),
		//newBug(g, bugsRed, screenWidth/2-10, eScreenHeight-100, bugDestroyFn),
		//newBug(g, bugsRed, screenWidth/2+10, eScreenHeight-100, bugDestroyFn),
		//newBug(g, bugsRed, screenWidth/2+30, eScreenHeight-100, bugDestroyFn),
		//newBug(g, bugsRed, screenWidth/2+50, eScreenHeight-100, bugDestroyFn),
	}

	for _, redBug := range redBugs {
		g.drawHandler.Add(redBug)
		g.updateHandler.Add(redBug)
		g.clickHandler.Add(redBug)
		g.AddEnemy(redBug)
	}

	//g.drawHandler.Add(newBug(g, bugsBlue, screenWidth/2, screenHeight-100))
	//g.drawHandler.Add(newBug(g, bugsGreen, screenWidth/2+50, screenHeight-100))

	// バリケードを家のすぐ下に配置
	//	barricadeOnDestroyFn := func(b *barricade) {
	//		g.drawHandler.Remove(b)
	//		g.clickHandler.Remove(b)
	//		g.RemoveBuilding(b)
	//		g.infoPanel.Remove(b)
	//	}

	barricades := []*barricade{
		//newBarricade(g, screenWidth/2-105, eScreenHeight/2+80, barricadeOnDestroyFn),
		//newBarricade(g, screenWidth/2, eScreenHeight/2+80, barricadeOnDestroyFn),
		//newBarricade(g, screenWidth/2+105, eScreenHeight/2+80, barricadeOnDestroyFn),
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

	// TODO: 本当はウェーブ中だけこれをやる
	//attackPane := newAttackPane(g)
	//g.clickHandler.Add(attackPane)

	// TODO: 本当はウェーブの間だけこれをやる
	buildPane := newBuildPane(g)
	g.clickHandler.Add(buildPane)
	g.drawHandler.Add(buildPane)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
