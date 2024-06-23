package main

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	_ "embed"
	_ "image/png"
)

//go:embed assets/house.png
var houseImageData []byte

type house struct {
	game *Game

	x, y          int // 画面中央に配置するので初期化時に値をもらう必要はない
	width, height int // 画像サイズをそのまま使うので初期化時に値をもらう必要はない
	zindex        int // これも適当に調整するので初期化時に値をもらう必要はない
	image         *ebiten.Image

	health int

	// 画像の拡大率。
	// TODO: 本当は画像のサイズそのものを変更したほうが見た目も処理効率も良くなる。余裕があれば後々やろう。
	scale float64

	// health が 0 になったときに呼ばれる関数
	onDestroy func(h *house)
}

func newHouse(game *Game) *house {
	img, _, err := image.Decode(bytes.NewReader(houseImageData))
	if err != nil {
		log.Fatal(err)
	}

	h := &house{
		game: game,

		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		scale:  0.5,

		health: 100,

		image: ebiten.NewImageFromImage(img),
	}

	h.x = screenWidth / 2
	h.y = eScreenHeight / 2

	return h
}

// 画面中央に配置
func (h *house) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(h.scale, h.scale)
	opts.GeoM.Translate(float64(h.x)-float64(h.width)*h.scale/2, float64(h.y)-float64(h.height)*h.scale/2)
	screen.DrawImage(h.image, opts)
}

func (h *house) ZIndex() int {
	return h.zindex
}

func (h *house) Position() (int, int) {
	// 中央の座標を返す
	return h.x, h.y
}

func (h *house) SetPosition(x, y int) {
	h.x = x
	h.y = y
}

func (h *house) Size() (int, int) {
	return int(float64(h.width) * h.scale), int(float64(h.height) * h.scale)
}

func (h *house) Name() string {
	return "House"
}

func (h *house) Damage(d int) {
	if h.health <= 0 {
		return
	}

	h.health -= d
	if h.health <= 0 {
		h.health = 0
	}
}

// house implements Clickable interface
func (h *house) OnClick(x, y int) bool {
	h.game.clickedObject = "House"
	// infoPanel に情報を表示する
	h.game.infoPanel.ClearButtons()
	icon := newHouseIcon(80, eScreenHeight+70)
	h.game.infoPanel.setIcon(icon)
	h.game.infoPanel.setUnit(h)

	switch h.game.phase {
	case PhaseBuilding:
		// infoPanel にバリケード建築ボタンを表示
		buildBarricadeButton := newButton(h.game,
			225, eScreenHeight, infoPanelHeight, infoPanelHeight, 1,
			func(x, y int) bool {
				barricadeOnDestroyFn := func(b *barricade) {
					b.game.drawHandler.Remove(b)
					b.game.clickHandler.Remove(b)
					b.game.RemoveBuilding(b)
					b.game.infoPanel.Remove(b)
				}

				h.game.buildCandidate = newBarricade(h.game, 0, 0, barricadeOnDestroyFn)
				return false
			},
			func(screen *ebiten.Image, x, y, width, height int) {
				drawRect(screen, x, y, width, height)
				barricadeIcon := newBarricadeIcon(x+width/2, y+height/2-10)
				barricadeIcon.Draw(screen)

				ebitenutil.DebugPrintAt(screen, "BUILD", x+width/2-20, y+height/2+40)
			})

		h.game.infoPanel.AddButton(buildBarricadeButton)
	case PhaseWave:
		// TODO: implement
	default:
		log.Fatalf("unexpected phase: %v", h.game.phase)
	}

	return false
}

func (h *house) Health() int {
	return h.health
}

// グレーアウトした drawRect を描画
func drawGrayRect(screen *ebiten.Image, x, y, width, height int) {
	strokeWidth := float32(2)
	vector.StrokeLine(screen, float32(x), float32(y), float32(x+width), float32(y), strokeWidth, color.Gray16{0x8888}, true)
	vector.StrokeLine(screen, float32(x), float32(y), float32(x), float32(y+height), strokeWidth, color.Gray16{0x8888}, true)
	vector.StrokeLine(screen, float32(x+width), float32(y), float32(x+width), float32(y+height), strokeWidth, color.Gray16{0x8888}, true)
	vector.StrokeLine(screen, float32(x), float32(y+height), float32(x+width), float32(y+height), strokeWidth, color.Gray16{0x8888}, true)
}

func drawRect(screen *ebiten.Image, x, y, width, height int) {
	strokeWidth := float32(2)
	vector.StrokeLine(screen, float32(x), float32(y), float32(x+width), float32(y), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(x), float32(y), float32(x), float32(y+height), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(x+width), float32(y), float32(x+width), float32(y+height), strokeWidth, color.White, true)
	vector.StrokeLine(screen, float32(x), float32(y+height), float32(x+width), float32(y+height), strokeWidth, color.White, true)
}

func (h *house) IsClicked(x, y int) bool {
	width, height := h.Size()
	return h.x-width/2 <= x && x <= h.x+width/2 && h.y-height/2 <= y && y <= h.y+height/2
}

func (h *house) SetOverlap(overlap bool) {
	// 登場の機会はないので実装しない
}

func (h *house) IsOverlap() bool {
	// 登場の機会はないので実装しない
	return false
}
