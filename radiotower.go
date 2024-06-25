package main

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
	_ "image/png"
)

//go:embed assets/radio_tower.png
var radioTowerImageData []byte

type radioTower struct {
	game *Game

	x, y          int
	width, height int
	zindex        int
	image         *ebiten.Image

	health           int
	shortAttackRange float64
	longAttackRange  float64
	attackZoneRadius float64
	attackPower      int
	cooldown         int
	erapsedTime      int // 攻撃実行からの経過時間

	// 画像の拡大率。
	// 1以外を指定する場合は元画像のサイズをそもそも変更できないか検討すること
	scale float64

	// health が 0 になったときに呼ばれる関数
	onDestroy func(b *radioTower)

	// この建物が他の建物と重なっているかどうか (建築確定前に用いるフラグ)
	isOverlapping bool
}

const radioTowerAttackCoolDown = 15

func newRadioTower(game *Game, x, y int, onDestroy func(b *radioTower)) *radioTower {
	img, _, err := image.Decode(bytes.NewReader(radioTowerImageData))
	if err != nil {
		log.Fatal(err)
	}

	h := &radioTower{
		game: game,

		x:      x,
		y:      y,
		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		scale:  1,

		health: 100,

		// 近すぎる敵は攻撃できない
		// 最長攻撃可能距離と、最短攻撃可能距離を設定する
		shortAttackRange: 200,
		longAttackRange:  400,
		attackZoneRadius: 50,

		attackPower: 10,

		image: ebiten.NewImageFromImage(img),

		onDestroy: onDestroy,
	}

	return h
}

func (t *radioTower) Update() {
	if t.game.phase == PhaseBuilding {
		// do nothing
		return
	}

	// 敵が攻撃範囲に入ってきたら攻撃する
	// 複数の敵が攻撃範囲に入ってきた場合は、最も近い敵を攻撃する
	// ただし近すぎる敵には攻撃できない

	// shortAttackRange と longAttackRange の間にいる敵のうち、最も近い敵を探す
	var nearestEnemy Enemy
	nearestDistance := math.MaxFloat64
	for _, e := range t.game.enemies {
		ex, ey := e.Position()
		// 敵が shortAttackRange と longAttackRange の間にいるかどうかを判定する
		distance := math.Sqrt(math.Pow(float64(t.x-ex), 2) + math.Pow(float64(t.y-ey), 2))
		if t.shortAttackRange < distance && distance < t.longAttackRange {
			if distance < nearestDistance {
				nearestEnemy = e
				nearestDistance = distance
			}
		}
	}

	// クールダウンが明けていて、攻撃可能な敵がいる場合は攻撃する
	if t.cooldown <= 0 && nearestEnemy != nil {
		// nearestEnemy を中心に範囲攻撃を行う
		ex, ey := nearestEnemy.Position()
		for _, e := range t.game.enemies {
			ex2, ey2 := e.Position()
			distance := math.Sqrt(math.Pow(float64(ex-ex2), 2) + math.Pow(float64(ey-ey2), 2))
			if distance < t.attackZoneRadius {
				b := e.(Damager)
				b.Damage(t.attackPower)
			}
		}

		t.cooldown = radioTowerAttackCoolDown

		// エフェクトを描画する
		eff := newRadioTowerAttackEffect(t.game, t.x, t.y)
		t.game.drawHandler.Add(eff)
	}

	if t.cooldown > 0 {
		t.cooldown--
	}
}

// タワーから発射されるビームを描画するための構造体
type radioTowerAttackEffect struct {
	game *Game

	x, y  int
	width int

	// 何フレーム後に消えるか
	displayTime int
}

func newRadioTowerAttackEffect(game *Game, x, y int) *radioTowerAttackEffect {
	return &radioTowerAttackEffect{
		game:        game,
		x:           x,
		y:           y,
		width:       7,
		displayTime: 10,
	}
}

func (b *radioTowerAttackEffect) Draw(screen *ebiten.Image) {
	if b.displayTime >= 0 {
		// TODO: implement
		b.displayTime--
	}

	if b.displayTime <= 0 {
		b.game.drawHandler.Remove(b)
	}
}

func (b *radioTowerAttackEffect) ZIndex() int {
	return 110
}

// 画面中央に配置
func (b *radioTower) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(b.scale, b.scale)
	opts.GeoM.Translate(float64(b.x)-float64(b.width)*b.scale/2, float64(b.y)-float64(b.height)*b.scale/2)

	// 他の建物と重なっている場合は赤くする
	if b.isOverlapping {
		opts.ColorScale.Scale(1, 0, 0, 1)
	} else if b.game.buildCandidate == b {
		// 建築確定前は暗い色で建物を描画する
		opts.ColorScale.Scale(0.5, 0.5, 0.5, 1)
	}

	screen.DrawImage(b.image, opts)
}

func (b *radioTower) ZIndex() int {
	return b.zindex
}

func (b *radioTower) Position() (int, int) {
	return b.x, b.y
}

func (b *radioTower) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

func (b *radioTower) Size() (int, int) {
	return int(float64(b.width) * b.scale), int(float64(b.height) * b.scale)
}

func (b *radioTower) Name() string {
	return "RadioTower"
}

func (b *radioTower) Damage(d int) {
	if b.health <= 0 {
		return
	}

	b.health -= d
	if b.health <= 0 {
		b.health = 0
		b.onDestroy(b)
	}
}

// radioTower implements Clickable interface
func (b *radioTower) OnClick(x, y int) bool {
	b.game.clickedObject = "radioTower"

	// infoPanel に情報を表示する

	// TODO: ClearButtons は呼び出し側でやるんじゃなくて infoPanel 側のどっかでやるべきかな
	b.game.infoPanel.ClearButtons()
	icon := newRadioTowerIcon(80, eScreenHeight+70)
	b.game.infoPanel.setIcon(icon)
	b.game.infoPanel.setUnit(b)

	return false
}

func newRadioTowerIcon(x, y int) *icon {
	img, _, err := image.Decode(bytes.NewReader(radioTowerImageData))
	if err != nil {
		log.Fatal(err)
	}

	return newIcon(x, y, ebiten.NewImageFromImage(img))
}

func (b *radioTower) Health() int {
	return b.health
}

func (b *radioTower) IsClicked(x, y int) bool {
	w, h := b.Size()
	return b.x-w/2 <= x && x <= b.x+w/2 && b.y-h/2 <= y && y <= b.y+h/2
}

func (b *radioTower) SetOverlap(overlap bool) {
	b.isOverlapping = overlap
}

func (b *radioTower) IsOverlap() bool {
	// 他の建物と重なっているかどうかを判定する
	for _, building := range b.game.buildings {
		if building == b {
			continue
		}

		bx, by := building.Position()
		bw, bh := building.Size()

		if intersects(
			rect{b.x - b.width/2, b.y - b.height/2, b.width, b.height},
			rect{bx - bw/2, by - bh/2, bw, bh},
		) {
			return true
		}
	}

	return false
}

func (b *radioTower) Cost() int {
	return CostRadioTowerBuild
}
