package main

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type title struct {
	game *Game

	stopFrame int

	houseImg     *ebiten.Image
	barricadeImg *ebiten.Image
	redBugImg    *ebiten.Image
	blueBugImg   *ebiten.Image
	greenBugImg  *ebiten.Image
}

func newTitle(g *Game) *title {
	// いったん BGM 止める
	getAudioPlayer().stopBGM()

	houseImg, _, err := image.Decode(bytes.NewReader(houseImageData))
	if err != nil {
		log.Fatal(err)
	}

	barricadeImg, _, err := image.Decode(bytes.NewReader(barricadeImageData))
	if err != nil {
		log.Fatal(err)
	}

	bugsImgD, _, err := image.Decode(bytes.NewReader(bugsImageData))
	if err != nil {
		log.Fatal(err)
	}
	bugsImg := ebiten.NewImageFromImage(bugsImgD)

	return &title{
		game:         g,
		stopFrame:    120,
		houseImg:     ebiten.NewImageFromImage(houseImg),
		barricadeImg: ebiten.NewImageFromImage(barricadeImg),
		redBugImg:    bugsImg.SubImage(redBug()).(*ebiten.Image),
		blueBugImg:   bugsImg.SubImage(blueBug()).(*ebiten.Image),
		greenBugImg:  bugsImg.SubImage(greenBug()).(*ebiten.Image),
	}
}

func (t *title) OnClick(x, y int) bool {
	// タイタオルバックを非表示にする
	getAudioPlayer().play(soundShot)
	t.game.clickHandler.Remove(t)
	t.game.updateHandler.Add(t)
	return false
}

func (t *title) Update() {
	// 数病経過したらゲーム画面に遷移する
	t.stopFrame--
	if t.stopFrame <= 0 {
		t.game.drawHandler.Remove(t)
		t.game.updateHandler.Remove(t)
		getAudioPlayer().playBGM()
	}
}

func (t *title) IsClicked(x, y int) bool {
	return true
}

func (t *title) Draw(screen *ebiten.Image) {
	// やや濃い緑で塗りつぶす
	screen.Fill(color.RGBA{0x00, 0x45, 0x00, 0xff})

	// 家を描く
	{
		width, height := t.houseImg.Bounds().Dx(), t.houseImg.Bounds().Dy()
		op := &ebiten.DrawImageOptions{}
		scaleX, scaleY := 2.5, 2.5
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(screenWidth-float64(width)*scaleX)/2-300, float64(screenHeight-float64(height)*scaleY)/2-200)
		screen.DrawImage(t.houseImg, op)
	}

	// バリケードを 3 つほど描く
	width, height := t.barricadeImg.Bounds().Dx(), t.barricadeImg.Bounds().Dy()
	scaleX, scaleY := 2.5, 2.5
	var op *ebiten.DrawImageOptions

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(screenWidth-float64(width)*scaleX)/2+200, float64(screenHeight-float64(height)*scaleY)/2-100)
	screen.DrawImage(t.barricadeImg, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(screenWidth-float64(width)*scaleX)/2, float64(screenHeight-float64(height)*scaleY)/2)
	screen.DrawImage(t.barricadeImg, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(screenWidth-float64(width)*scaleX)/2-200, float64(screenHeight-float64(height)*scaleY)/2+100)
	screen.DrawImage(t.barricadeImg, op)

	// 虫を描く
	{
		scaleX, scaleY := float64(5), float64(5)

		width, height := t.redBugImg.Bounds().Dx(), t.redBugImg.Bounds().Dy()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(screenWidth-float64(width)*scaleX)/2+200+100, float64(screenHeight-float64(height)*scaleY)/2+100)
		screen.DrawImage(t.redBugImg, op)

		width, height = t.blueBugImg.Bounds().Dx(), t.blueBugImg.Bounds().Dy()
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(screenWidth-float64(width)*scaleX)/2+100, float64(screenHeight-float64(height)*scaleY)/2+200)
		screen.DrawImage(t.blueBugImg, op)

		width, height = t.greenBugImg.Bounds().Dx(), t.greenBugImg.Bounds().Dy()
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(screenWidth-float64(width)*scaleX)/2-200+100, float64(screenHeight-float64(height)*scaleY)/2+300)
		screen.DrawImage(t.greenBugImg, op)
	}

	// 文字を描く
	clr := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	scaleX, scaleY = float64(5), float64(5)
	drawText(screen, "HOUSE DEFENSE OPERATION!", screenWidth-750, 100, scaleX, scaleY, clr)
	drawText(screen, "CLICK TO START!", screenWidth-750, 170, scaleX, scaleY, clr)
}

func (t *title) ZIndex() int {
	return 300
}
