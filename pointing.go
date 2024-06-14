package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func getClickedPosition() (x, y int, eventOccurred bool) {
	// マウスクリックの処理
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		return mx, my, true
	}

	// タッチイベントの処理
	touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
	if len(touchIDs) > 0 {
		tx, ty := ebiten.TouchPosition(touchIDs[0])
		return tx, ty, true
	}

	// イベントが発生しなかった場合
	return 0, 0, false
}
