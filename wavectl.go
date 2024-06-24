package main

type waveController struct {
	game *Game

	onWaveEnd func()
}

func newWaveController(g *Game, onWaveEnd func()) *waveController {
	return &waveController{
		game:      g,
		onWaveEnd: onWaveEnd,
	}
}

func (w *waveController) Update() {
	// buildings から house を取得
	// TODO: いちいちループ回すのは効率が悪いかも

	if w.game.house.health <= 0 {
		// ゲームオーバーの処理
		gover := newGameover(w.game)
		w.game.clickHandler.Add(gover)
		w.game.drawHandler.Add(gover)

		// ウェーブ終了みたいなものなので自分を削除する
		w.game.updateHandler.Remove(w)
	}

	if len(w.game.enemies) == 0 {
		w.onWaveEnd()
		// ウェーブが終了したら自分自身を削除する
		w.game.updateHandler.Remove(w)
	}
}
