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
	// デバッグのための変数
	clickedPositionX, clickedPositionY int
	clickedObject                      string

	// ゲーム全編通して使うハンドラ
	clickHandler  *OnClickHandler
	drawHandler   *DrawHandler
	updateHandler *UpdateHandler

	// 以下はメインのゲームシーンで使う変数
	// TODO: シーンごとに Game 構造体を分けるべきかもしれない
	phase Phase

	// 建物のリスト
	buildings []Building

	// 敵のリスト
	enemies []Enemy

	// 情報パネル
	infoPanel *infoPanel

	// 建築対象としていったん保持されているオブジェクト
	buildCandidate Building

	// panes
	attackPane *attackPane
	buildPane  *buildPane

	// ウェーブのコントローラ
	waveCtrl *waveController

	credit int
}

type Phase int

const (
	screenWidth  = 1280
	screenHeight = 960
)

const (
	// 画面上にデバッグ情報を表示するかどうか
	debugEnabled = true
)

const (
	// 建築フェーズ
	PhaseBuilding Phase = iota
	// ウェーブフェーズ
	PhaseWave
)

// コスト一覧
const (
	CostBarricadeBuild = 30
)

const (
	// infoPanel の高さを計算
	// infoPanel の高さの分だけ、ゲーム画面の中央座標が上にずれる
	// 中央座標計算のためにあらかじめここで計算しておく
	infoPanelHeight = screenHeight / 7
	eScreenHeight   = screenHeight - infoPanelHeight - 10
)

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
	g.drawHandler.HandleDraw(screen)

	// 現在のフェーズを表示
	switch g.phase {
	case PhaseBuilding:
		ebitenutil.DebugPrintAt(screen, "Phase: Building", 0, 40)
	case PhaseWave:
		ebitenutil.DebugPrintAt(screen, "Phase: Wave", 0, 40)
	}

	// 画面右上にクレジットを表示
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Credit: %d", g.credit), screenWidth-100, 0)

	// 以下はデバッグ情報

	if debugEnabled {
		// クリックされた位置を表示
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Clicked Position: (%d, %d)", g.clickedPositionX, g.clickedPositionY), 0, 0)
		// クリックされたオブジェクトを表示
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Clicked Object: %s", g.clickedObject), 0, 20)

		// drawHandler の長さを表示
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DrawHandler: %d", len(g.drawHandler.drawable)), 0, 60)
		// clickHandler の長さを表示
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ClickHandler: %d", len(g.clickHandler.clickableObjects)), 0, 80)
		// updateHandler の長さを表示
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("UpdateHandler: %d", len(g.updateHandler.updaters)), 0, 100)

		// 画面中央に点を表示 (debug)
		vector.DrawFilledRect(screen, screenWidth/2, eScreenHeight/2, 1, 1, color.RGBA{255, 255, 255, 255}, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) SetBuildingPhase() {
	g.phase = PhaseBuilding

	// 情報パネルをいったんクリアする
	g.infoPanel.unit = nil
	g.infoPanel.ClearButtons()

	// wave phase で追加したものを削除
	g.attackPane.RemoveAll()

	// Build phase に必要なものを追加
	g.buildPane = newBuildPane(g)
	g.clickHandler.Add(g.buildPane)
	g.drawHandler.Add(g.buildPane)
}

func (g *Game) SetWavePhase() {
	g.phase = PhaseWave

	// 情報パネルをクリアする (Build のメニューなどが出ていたら消すため)
	g.infoPanel.unit = nil
	g.infoPanel.ClearButtons()

	// building phase で追加したものを削除
	g.buildPane.RemoveAll()

	// Wave phase に必要なものを追加
	g.attackPane = newAttackPane(g)
	g.clickHandler.Add(g.attackPane)

	// とりあえずいったん虫を画面の下部に配置
	// TODO: wave の設定にしたがって敵を生成できるようにする
	bugDestroyFn := func(b *bug) {
		g.drawHandler.Remove(b)
		g.updateHandler.Remove(b)
		g.clickHandler.Remove(b)
		g.RemoveEnemy(b)
		g.infoPanel.Remove(b)
	}

	//とりあえずいったん虫を画面の下部に配置
	redBugs := []*bug{
		newBug(g, bugsRed, screenWidth/2-50, eScreenHeight-100, bugDestroyFn),
		newBug(g, bugsRed, screenWidth/2-30, eScreenHeight-100, bugDestroyFn),
		newBug(g, bugsRed, screenWidth/2-10, eScreenHeight-100, bugDestroyFn),
		newBug(g, bugsRed, screenWidth/2+10, eScreenHeight-100, bugDestroyFn),
		newBug(g, bugsRed, screenWidth/2+30, eScreenHeight-100, bugDestroyFn),
		newBug(g, bugsRed, screenWidth/2+50, eScreenHeight-100, bugDestroyFn),
	}

	for _, redBug := range redBugs {
		g.drawHandler.Add(redBug)
		g.updateHandler.Add(redBug)
		g.clickHandler.Add(redBug)
		g.AddEnemy(redBug)
	}

	// 敵が全滅したらウェーブを終了して建築フェーズに戻る
	// 敵が全滅したことをコールバックする
	waveEndFn := func() {
		g.SetBuildingPhase()

		// ウェーブ終了時に一定のクレジットを得る
		g.credit += 100
	}

	g.waveCtrl = newWaveController(g, waveEndFn)
	g.updateHandler.Add(g.waveCtrl)

	//g.drawHandler.Add(newBug(g, bugsBlue, screenWidth/2, screenHeight-100))
	//g.drawHandler.Add(newBug(g, bugsGreen, screenWidth/2+50, screenHeight-100))
}

func (g *Game) initialize() {
	// とりあえずいきなりゲームが始まるとする。
	// TODO: まずタイトルバックを表示して、その後にゲーム画面に遷移するようにする
	house := newHouse(g)
	g.drawHandler.Add(house)
	g.AddBuilding(house)
	g.clickHandler.Add(house)

	g.infoPanel = newInfoPanel(g, screenWidth-20, infoPanelHeight)
	g.drawHandler.Add(g.infoPanel)

	// 背景担当
	bg := newBackground(g)
	g.drawHandler.Add(bg)

	// クレジットを初期化
	g.credit = 100
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")

	g := &Game{
		clickHandler:  &OnClickHandler{},
		drawHandler:   &DrawHandler{},
		updateHandler: &UpdateHandler{},
	}

	g.initialize()

	// 最初のシーンをセットアップする
	g.SetBuildingPhase()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
