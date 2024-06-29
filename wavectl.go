package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type waveController struct {
	game *Game

	currentBigWave   int
	currentSmallWave int

	erapsedFrame int
	onWaveEnd    func()
}

func newWaveController(g *Game, onWaveEnd func()) *waveController {
	return &waveController{
		game: g,

		currentBigWave:   0,
		currentSmallWave: 0,
		erapsedFrame:     0,
		onWaveEnd:        onWaveEnd,
	}
}

func (w *waveController) Update() {
	w.spawnEnemy()
	w.erapsedFrame++

	if w.game.house.health <= 0 {
		// ゲームオーバーの処理
		gover := newGameover(w.game)
		w.game.updateHandler.Add(gover)
		w.game.clickHandler.Add(gover)
		w.game.drawHandler.Add(gover)

		getAudioPlayer().stopBGM()
		getAudioPlayer().play(soundGameover)

		// ウェーブ終了みたいなものなので自分を削除する
		w.game.updateHandler.Remove(w)
	} else if len(w.game.enemies) == 0 {
		// enemies が 0 になるということは、small wave が終わったか、big wave が終わったということ
		// TIPS: なので、ウェーブが始まったら最初のフレームでかならず enemies を 1 以上にすること。
		// そうでないとウェーブがはじまった瞬間にウェーブが終わってしまう

		w.onWaveEnd()
		w.erapsedFrame = 0
		w.currentSmallWave = 0
		w.currentBigWave++

		// ゲームクリアの処理
		if w.currentBigWave == len(waveList) {
			getAudioPlayer().stopBGM()
			getAudioPlayer().play(soundClear)

			gclear := newGameClear(w.game)
			w.game.updateHandler.Add(gclear)
			w.game.clickHandler.Add(gclear)
			w.game.drawHandler.Add(gclear)
		} else {
			// ウェーブ間の処理
			t := newTimerText(w.game, screenWidth/2-350, screenHeight/2+50, "Wave Clear! Credit Earned! $120")
			w.game.drawHandler.Add(t)
			w.game.updateHandler.Add(t)
			t = newTimerText(w.game, screenWidth/2-200, screenHeight/2+150, fmt.Sprintf("Waves remaining: %d", len(waveList)-w.currentBigWave))
			w.game.drawHandler.Add(t)
			w.game.updateHandler.Add(t)

		}

		// ウェーブが終了したら自分自身を削除する
		w.game.updateHandler.Remove(w)
	}
}

// ウェーブ間に表示するテキスト
// 主にお金が手に入ったことを伝えるのが目的
type timerText struct {
	game *Game

	x, y int
	text string

	displayFrame int
}

func newTimerText(g *Game, x, y int, text string) *timerText {
	return &timerText{
		game:         g,
		x:            x,
		y:            y,
		text:         text,
		displayFrame: 300,
	}
}

func (c *timerText) Update() {
	c.displayFrame--
	if c.displayFrame <= 0 {
		c.game.drawHandler.Remove(c)
		c.game.updateHandler.Remove(c)
	}
}

func (c *timerText) Draw(screen *ebiten.Image) {
	drawText(screen, c.text, c.x, c.y, 4, 4, color.RGBA{0xff, 0xff, 0xff, 0xff})
}

func (c *timerText) ZIndex() int {
	return 300
}

type spawnInfo struct {
	color bugColor
	x, y  int
}

// トータル10になるようにする
type bugSpawnRatio struct {
	red, blue, green int
}

func generateSpawnInfos(num int, spawnRatio bugSpawnRatio) []spawnInfo {
	rand.NewSource(time.Now().UnixNano())
	var infos []spawnInfo

	for i := 0; i < num; i++ {
		// 四方八方からランダムに生成する
		side := rand.Intn(4) // 0: 上, 1: 下, 2: 左, 3: 右
		var x, y int
		switch side {
		case 0: // 上
			x = rand.Intn(screenWidth)
			y = -50
		case 1: // 下
			x = rand.Intn(screenWidth)
			y = screenHeight + 50
		case 2: // 左
			x = -50
			y = rand.Intn(screenHeight)
		case 3: // 右
			x = screenWidth + 50
			y = rand.Intn(screenHeight)
		}

		r := rand.Intn(10)
		if r < spawnRatio.red {
			infos = append(infos, spawnInfo{bugsRed, x, y})
		} else if r < spawnRatio.red+spawnRatio.blue {
			infos = append(infos, spawnInfo{bugsBlue, x, y})
		} else {
			infos = append(infos, spawnInfo{bugsGreen, x, y})
		}
	}

	return infos
}

var waveList = [][]struct {
	spawnFrame    int
	spawnInfoList []spawnInfo
}{
	// ウェーブにおける敵の戦闘力は以下のように計算してみる
	// 1. 赤虫: 1, 青虫: 2, 緑虫: 3
	// 2. それぞれの虫の数をかけて、それを足し合わせる
	// 3. それをウェーブの戦闘力とする
	// 例: 赤虫が 5, 青虫が 3, 緑虫が 2 の場合、戦闘力は 5*1 + 3*2 + 2*3 = 5 + 6 + 6 = 17 となる
	// 後半のウェーブは戦闘力が高くなるように設定している

	{ // 戦闘力10 赤だけ
		{0, generateSpawnInfos(5, bugSpawnRatio{10, 0, 0})},
		{60, generateSpawnInfos(5, bugSpawnRatio{10, 0, 0})},
	},
	{ // 戦闘力20 青だけ
		{0, generateSpawnInfos(5, bugSpawnRatio{0, 10, 0})},
		{60, generateSpawnInfos(5, bugSpawnRatio{0, 10, 0})},
	},
	{ // 戦闘力30 緑だけ
		{0, generateSpawnInfos(5, bugSpawnRatio{0, 0, 10})},
		{60, generateSpawnInfos(5, bugSpawnRatio{0, 0, 10})},
	},
	{ // 戦闘力40 赤青混合
		{0, generateSpawnInfos(12, bugSpawnRatio{4, 6, 0})},
		{60, generateSpawnInfos(13, bugSpawnRatio{4, 6, 0})},
	},
	{ // 戦闘力50 青緑混合
		{0, generateSpawnInfos(12, bugSpawnRatio{0, 6, 4})},
		{60, generateSpawnInfos(13, bugSpawnRatio{0, 6, 4})},
	},
	{ // 戦闘力60 赤緑混合
		{0, generateSpawnInfos(19, bugSpawnRatio{7, 0, 3})},
		{60, generateSpawnInfos(19, bugSpawnRatio{7, 0, 3})},
	},
	{ // 戦闘力70 全部混合ちょっといっぱいくる
		{0, generateSpawnInfos(20, bugSpawnRatio{3, 5, 2})},
		{60, generateSpawnInfos(20, bugSpawnRatio{3, 5, 2})},
		{120, generateSpawnInfos(20, bugSpawnRatio{3, 5, 2})},
		{240, generateSpawnInfos(20, bugSpawnRatio{3, 5, 2})},
	},
	{ // 戦闘力80 全部混合ちょっと控えめ
		{0, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
		{60, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
		{120, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
	},
	{ // 戦闘力90 全部混合ちょっと控えめ
		{0, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
		{60, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
		{120, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
	},
	{ // 戦闘力90 全部混合ちょっと控えめ
		{0, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
		{60, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
		{120, generateSpawnInfos(14, bugSpawnRatio{3, 5, 2})},
	},
	{ // 戦闘力100 全部混合いっぱいくる
		{0, generateSpawnInfos(30, bugSpawnRatio{3, 5, 2})},
		{60, generateSpawnInfos(30, bugSpawnRatio{3, 5, 2})},
		{120, generateSpawnInfos(30, bugSpawnRatio{3, 5, 2})},
		{240, generateSpawnInfos(30, bugSpawnRatio{3, 5, 2})},
	},
}

func (w *waveController) spawnEnemy() []Enemy {
	// erapsedFrame に従って敵を生成する

	bugDestroyFn := func(b *bug) {
		w.game.drawHandler.Remove(b)
		w.game.updateHandler.Remove(b)
		w.game.clickHandler.Remove(b)
		w.game.RemoveEnemy(b)
		w.game.infoPanel.Remove(b)
	}

	if w.currentBigWave >= len(waveList) {
		return nil
	}
	if w.currentSmallWave >= len(waveList[w.currentBigWave]) {
		return nil
	}

	if w.erapsedFrame == waveList[w.currentBigWave][w.currentSmallWave].spawnFrame {
		spawnList := waveList[w.currentBigWave][w.currentSmallWave].spawnInfoList
		enemies := make([]Enemy, 0, len(spawnList))
		for _, info := range spawnList {
			enemies = append(enemies, newBug(w.game, info.color, info.x, info.y, bugDestroyFn))
		}

		for _, e := range enemies {
			w.game.drawHandler.Add(e)
			w.game.updateHandler.Add(e)
			w.game.clickHandler.Add(e)
			w.game.AddEnemy(e)
		}
		w.currentSmallWave++

		return enemies
	}

	return nil
}
