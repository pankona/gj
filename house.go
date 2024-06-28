package main

import (
	"bytes"
	"fmt"
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

	// 死亡時のアニメーションを管理するための変数
	deadAnimationDuration int

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

		onDestroy: func(h *house) {
			// TODO: 爆発したり消えたりする処理を書く
			// ここでフラグを設定しといて、Update() や Draw で続きの処理を行うのもあり
			// いったんシンプルに消す
			h.game.updateHandler.Remove(h)
			h.game.RemoveBuilding(h)
			h.game.drawHandler.Remove(h)
		},
	}

	h.x = screenWidth / 2
	h.y = eScreenHeight / 2

	return h
}

func (h *house) Update() {
	if h.health <= 0 {
		h.deadAnimationDuration++
		if h.deadAnimationDuration >= deadAnimationTotalFrame {
			h.onDestroy(h)
		}
		return
	}
}

// 画面中央に配置
func (h *house) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	if h.health <= 0 {
		// 死亡時のアニメーションを行う
		// ぺちゃんこになるように縮小する
		// TODO: ちょっとアニメーションが怪しいので調整する
		scale := h.scale * (1.0 - float64(h.deadAnimationDuration)/deadAnimationTotalFrame)
		if scale < 0 {
			scale = 0
		}

		opts.GeoM.Translate(0, h.scale*float64(-h.height)/2)
		opts.GeoM.Scale(h.scale, scale)
		opts.GeoM.Translate(0, h.scale*float64(h.height)/2)
	} else {
		opts.GeoM.Scale(h.scale, h.scale)
	}

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

	// 建築 instruction を消す
	if h.game.buildInstruction != nil {
		h.game.drawHandler.Remove(h.game.buildInstruction)
		h.game.buildInstruction = nil
	}

	getAudioPlayer().play(soundChoice)

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
				if h.game.credit < CostBarricadeBuild {
					// お金が足りない場合は建築できない
					return false
				}

				getAudioPlayer().play(soundChoice)

				// buildCandidate を持っているときにバリケードボタンを押したときの振る舞い
				// 選択肢なおしということ、いったん手放す
				if h.game.buildCandidate != nil {
					h.game.drawHandler.Remove(h.game.buildCandidate)
				}

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

				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BUILD ($%d)", CostBarricadeBuild), x+width/2-30, y+height/2+40)

				// 選択中であればボタンをハイライト表示する
				if h.game.buildCandidate != nil && h.game.buildCandidate.Name() == "Barricade" {
					drawYellowRect(screen, x, y, width, height)
				}

				// お金が足りないときはボタン全体をグレーアウトする
				if h.game.credit < CostBarricadeBuild {
					overlay := ebiten.NewImage(width, height)
					overlay.Fill(color.RGBA{128, 128, 128, 128})
					overlayOpts := &ebiten.DrawImageOptions{}
					overlayOpts.GeoM.Translate(float64(x), float64(y))
					screen.DrawImage(overlay, overlayOpts)
				}
			})

		h.game.infoPanel.AddButton(buildBarricadeButton)

		buildTowerButton := newButton(h.game,
			225+infoPanelHeight, eScreenHeight, infoPanelHeight, infoPanelHeight, 1,
			func(x, y int) bool {
				if h.game.credit < CostTowerBuild {
					// お金が足りない場合は建築できない
					return false
				}

				getAudioPlayer().play(soundChoice)

				// buildCandidate を持っているときにバリケードボタンを押したときの振る舞い
				// 選択肢なおしということ、いったん手放す
				if h.game.buildCandidate != nil {
					h.game.drawHandler.Remove(h.game.buildCandidate)
				}

				towerOnDestroyFn := func(b *tower) {
					b.game.updateHandler.Remove(b)
					b.game.drawHandler.Remove(b)
					b.game.clickHandler.Remove(b)
					b.game.RemoveBuilding(b)
					b.game.infoPanel.Remove(b)
				}

				h.game.buildCandidate = newTower(h.game, 0, 0, towerOnDestroyFn)
				return false
			},
			func(screen *ebiten.Image, x, y, width, height int) {
				drawRect(screen, x, y, width, height)
				towerIcon := newTowerIcon(x+width/2, y+height/2-10)
				towerIcon.Draw(screen)

				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BUILD ($%d)", CostTowerBuild), x+width/2-30, y+height/2+40)

				// 選択中であればボタンをハイライト表示する
				if h.game.buildCandidate != nil && h.game.buildCandidate.Name() == "Tower" {
					drawYellowRect(screen, x, y, width, height)
				}

				// お金が足りないときはボタン全体をグレーアウトする
				if h.game.credit < CostTowerBuild {
					overlay := ebiten.NewImage(width, height)
					overlay.Fill(color.RGBA{128, 128, 128, 128})
					overlayOpts := &ebiten.DrawImageOptions{}
					overlayOpts.GeoM.Translate(float64(x), float64(y))
					screen.DrawImage(overlay, overlayOpts)
				}
			})
		h.game.infoPanel.AddButton(buildTowerButton)

		buildRadioTowerButton := newButton(h.game,
			225+infoPanelHeight*2, eScreenHeight, infoPanelHeight, infoPanelHeight, 1,
			func(x, y int) bool {
				if h.game.credit < CostRadioTowerBuild {
					// お金が足りない場合は建築できない
					return false
				}

				getAudioPlayer().play(soundChoice)

				// buildCandidate を持っているときにバリケードボタンを押したときの振る舞い
				// 選択肢なおしということ、いったん手放す
				if h.game.buildCandidate != nil {
					h.game.drawHandler.Remove(h.game.buildCandidate)
				}

				radioTowerOnDestroyFn := func(b *radioTower) {
					b.game.updateHandler.Remove(b)
					b.game.drawHandler.Remove(b)
					b.game.clickHandler.Remove(b)
					b.game.RemoveBuilding(b)
					b.game.infoPanel.Remove(b)
				}

				h.game.buildCandidate = newRadioTower(h.game, 0, 0, radioTowerOnDestroyFn)
				return false
			},
			func(screen *ebiten.Image, x, y, width, height int) {
				drawRect(screen, x, y, width, height)
				radioTowerIcon := newRadioTowerIcon(x+width/2, y+height/2-10)
				radioTowerIcon.Draw(screen)

				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BUILD ($%d)", CostRadioTowerBuild), x+width/2-30, y+height/2+40)

				// 選択中であればボタンをハイライト表示する
				if h.game.buildCandidate != nil && h.game.buildCandidate.Name() == "RadioTower" {
					drawYellowRect(screen, x, y, width, height)
				}

				// お金が足りないときはボタン全体をグレーアウトする
				if h.game.credit < CostRadioTowerBuild {
					overlay := ebiten.NewImage(width, height)
					overlay.Fill(color.RGBA{128, 128, 128, 128})
					overlayOpts := &ebiten.DrawImageOptions{}
					overlayOpts.GeoM.Translate(float64(x), float64(y))
					screen.DrawImage(overlay, overlayOpts)
				}
			})
		h.game.infoPanel.AddButton(buildRadioTowerButton)

		nextWaveStartButton := newButton(h.game,
			screenWidth-10-infoPanelHeight, eScreenHeight, infoPanelHeight, infoPanelHeight, 1,
			func(x, y int) bool {
				getAudioPlayer().play(soundKettei)

				switch h.game.phase {
				case PhaseBuilding:
					// 最初のウェーブだったら攻撃方法に関する説明を表示する
					if h.game.waveCtrl.currentBigWave == 0 {
						h.game.attackInstruction = newInstruction(h.game, "CLICK BUGS TO ATTACK!", screenWidth/2-60, eScreenHeight/2+50)
						h.game.drawHandler.Add(h.game.attackInstruction)
					}
					h.game.SetWavePhase()
				case PhaseWave:
					// never reach
					fallthrough
				default:
					log.Fatalf("unexpected phase: %v", h.game.phase)
				}

				return false
			},
			func(screen *ebiten.Image, x, y, width, height int) {
				drawRect(screen, x, y, width, height)
				ebitenutil.DebugPrintAt(screen, "FINISH BUILDING!", x+width/2-45, y+height/2-40)
				ebitenutil.DebugPrintAt(screen, "START NEXT WAVE!", x+width/2-45, y+height/2-8)
				// 現在のウェーブとトータルウェーブ数を表示する
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CURRENT WAVE: %d/%d", h.game.waveCtrl.currentBigWave, len(waveList)), x+width/2-52, y+height/2+32)
			},
		)
		h.game.infoPanel.AddButton(nextWaveStartButton)

	case PhaseWave:
		// TODO: implement
	default:
		log.Fatalf("unexpected phase: %v", h.game.phase)
	}

	return false
}

func drawYellowRect(screen *ebiten.Image, x, y, width, height int) {
	strokeWidth := float32(10)
	// 黄色っぽい線で描画
	vector.StrokeLine(screen, float32(x), float32(y+5), float32(x+width), float32(y+5), strokeWidth, color.RGBA{0xff, 0xff, 0x00, 0xff}, true)
	vector.StrokeLine(screen, float32(x+5), float32(y+5), float32(x+5), float32(y+height-5), strokeWidth, color.RGBA{0xff, 0xff, 0x00, 0xff}, true)
	vector.StrokeLine(screen, float32(x+width-5), float32(y+5), float32(x+width-5), float32(y+height-5), strokeWidth, color.RGBA{0xff, 0xff, 0x00, 0xff}, true)
	vector.StrokeLine(screen, float32(x), float32(y+height-5), float32(x+width), float32(y+height-5), strokeWidth, color.RGBA{0xff, 0xff, 0x00, 0xff}, true)
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

func (h *house) Cost() int {
	// 登場の機会はないので実装しない
	return 0
}
