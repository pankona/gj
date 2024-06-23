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
	if len(w.game.enemies) == 0 {
		w.onWaveEnd()
		// ウェーブが終了したら自分自身を削除する
		w.game.updateHandler.Remove(w)
	}
}
