package main

import (
	"math/rand"
	"time"
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
		w.game.clickHandler.Add(gover)
		w.game.drawHandler.Add(gover)

		// 建物の Update() が呼ばれないようにする
		// ゲームオーバー画面で攻撃を続けないようにするため
		for _, b := range w.game.buildings {
			w.game.updateHandler.Remove(b)
		}

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
			gclear := newGameClear(w.game)
			w.game.clickHandler.Add(gclear)
			w.game.drawHandler.Add(gclear)

		} else {
			// ウェーブ間の処理
			// TODO: 必要ならなにか実装する
		}

		// ウェーブが終了したら自分自身を削除する
		w.game.updateHandler.Remove(w)
	}
}

type spawnInfo struct {
	color bugColor
	x, y  int
}

func generateSpawnInfos(num int) []spawnInfo {
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

		infos = append(infos, spawnInfo{
			color: bugsBlue,
			x:     x,
			y:     y,
		})
	}

	return infos
}

var waveList = [][]struct {
	spawnFrame    int
	spawnInfoList []spawnInfo
}{
	{
		{0, generateSpawnInfos(20)},
		{60, generateSpawnInfos(20)},
		{120, generateSpawnInfos(20)},
	},
	{
		{0, generateSpawnInfos(25)},
		{60, generateSpawnInfos(25)},
		{120, generateSpawnInfos(25)},
	},
}

func (w *waveController) spawnEnemy() []Enemy {
	// erapsedFrame に従って敵を生成する

	// いずれ青虫や緑虫も出すようにする
	//g.drawHandler.Add(newBug(g, bugsBlue, screenWidth/2, screenHeight-100))
	//g.drawHandler.Add(newBug(g, bugsGreen, screenWidth/2+50, screenHeight-100))

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
