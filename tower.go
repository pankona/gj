package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	_ "embed"
	"image/color"
	_ "image/png"
)

//go:embed assets/tower.png
var towerImageData []byte

type tower struct {
	game *Game

	x, y          int
	width, height int
	zindex        int
	image         *ebiten.Image

	health      int
	attackRange float64
	attackPower int
	cooldown    int
	erapsedTime int // 攻撃実行からの経過時間

	// 画像の拡大率。
	// 1以外を指定する場合は元画像のサイズをそもそも変更できないか検討すること
	scale float64

	// 死亡時のアニメーションを管理するための変数
	deadAnimationDuration int

	// health が 0 になったときに呼ばれる関数
	onDestroy func(b *tower)

	// この建物が他の建物と重なっているかどうか (建築確定前に用いるフラグ)
	isOverlapping bool
}

const towerAttackCoolDown = 30

func newTower(game *Game, x, y int, onDestroy func(b *tower)) *tower {
	img, _, err := image.Decode(bytes.NewReader(towerImageData))
	if err != nil {
		log.Fatal(err)
	}

	h := &tower{
		game: game,

		x:      x,
		y:      y,
		width:  img.Bounds().Dx(),
		height: img.Bounds().Dy(),
		scale:  1,

		health:      70,
		attackRange: 300,
		attackPower: 1,

		image: ebiten.NewImageFromImage(img),

		onDestroy: onDestroy,
	}

	return h
}

func (t *tower) Update() {
	if t.game.phase == PhaseBuilding {
		// do nothing
		return
	}

	// 死亡時のアニメーションを再生する
	if t.health <= 0 {
		t.deadAnimationDuration++
		if t.deadAnimationDuration >= deadAnimationTotalFrame {
			t.onDestroy(t)
		}
		return
	}

	// 家が壊れていたらもはや攻撃をやめる
	if t.game.house.health <= 0 {
		return
	}

	// 敵が攻撃範囲に入ってきたら攻撃する
	// 複数の敵が攻撃範囲に入ってきた場合は、最も近い敵を攻撃する

	// 最寄りの敵を探す
	var nearestEnemy Enemy
	nearestDistance := math.MaxFloat64
	for _, e := range t.game.enemies {
		ex, ey := e.Position()
		distance := math.Sqrt(math.Pow(float64(t.x-ex), 2) + math.Pow(float64(t.y-ey), 2))
		if distance < nearestDistance {
			nearestEnemy = e
			nearestDistance = distance
		}
	}

	// クールダウンが明けていて、かつ攻撃範囲に入っていれば攻撃する
	if t.cooldown == 0 && nearestEnemy != nil && nearestDistance < t.attackRange {
		bx, by := nearestEnemy.Position()
		b := nearestEnemy.(Damager)

		getAudioPlayer().play(soundBeam)

		b.Damage(t.attackPower)
		t.cooldown = towerAttackCoolDown

		// ビームを描画する
		bm := newBeam(t.game, t.x, t.y, bx, by)
		t.game.drawHandler.Add(bm)
	}

	if t.cooldown > 0 {
		t.cooldown--
	}
}

// タワーから発射されるビームを描画するための構造体
type beam struct {
	game *Game

	startX, startY int
	endX, endY     int
	width          int

	// 何フレーム後に消えるか
	displayTime int
}

func newBeam(game *Game, startX, startY, endX, endY int) *beam {
	return &beam{
		game:        game,
		startX:      startX,
		startY:      startY,
		endX:        endX,
		endY:        endY,
		width:       7,
		displayTime: 10,
	}
}

func (b *beam) Draw(screen *ebiten.Image) {
	if b.displayTime >= 0 {
		vector.StrokeLine(screen, float32(b.startX), float32(b.startY), float32(b.endX), float32(b.endY), float32(b.width), color.RGBA{255, 255, 0, 128}, true)
		b.displayTime--
	}

	if b.displayTime <= 0 {
		b.game.drawHandler.Remove(b)
	}
}

func (b *beam) ZIndex() int {
	return 110
}

// 画面中央に配置
func (b *tower) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}

	if b.health <= 0 {
		// 死亡時のアニメーションを行う
		// ぺちゃんこになるように縮小する
		scale := 1.0 - float64(b.deadAnimationDuration)/deadAnimationTotalFrame
		if scale < 0 {
			scale = 0
		}

		opts.GeoM.Translate(0, float64(-b.height))
		opts.GeoM.Scale(1, scale)
		opts.GeoM.Translate(0, float64(b.height))
	} else {
		opts.GeoM.Scale(b.scale, b.scale)
	}

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

func (b *tower) ZIndex() int {
	return b.zindex
}

func (b *tower) Position() (int, int) {
	return b.x, b.y
}

func (b *tower) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

func (b *tower) Size() (int, int) {
	return int(float64(b.width) * b.scale), int(float64(b.height) * b.scale)
}

func (b *tower) Name() string {
	return "Tower"
}

func (b *tower) Damage(d int) {
	if b.health <= 0 {
		return
	}

	b.health -= d
	if b.health <= 0 {
		getAudioPlayer().play(soundKuzureru)
		b.health = 0
	}
}

// tower implements Clickable interface
func (b *tower) OnClick(x, y int) bool {
	if b.game.buildCandidate != nil {
		// 建築予定のものを持っているときには何もしない
		return false
	}

	b.game.clickedObject = "tower"
	getAudioPlayer().play(soundChoice)

	// infoPanel に情報を表示する

	// TODO: ClearButtons は呼び出し側でやるんじゃなくて infoPanel 側のどっかでやるべきかな
	b.game.infoPanel.ClearButtons()
	icon := newTowerIcon(80, eScreenHeight+70)
	b.game.infoPanel.setIcon(icon)
	b.game.infoPanel.setUnit(b)
	b.game.infoPanel.drawDescriptionFn = func(screen *ebiten.Image, x, y int) {
		var scale float64 = 2
		// 敵を一匹ずつ攻撃するという説明を記載する
		drawText(screen, "I am Beam Tower!", x, y-10, scale, scale, color.RGBA{0xff, 0xff, 0xff, 0xff})
		drawText(screen, fmt.Sprintf("Cost: $%d", b.Cost()), x, y+20, scale, scale, color.RGBA{0xff, 0xff, 0xff, 0xff})
		drawText(screen, "Attack single bug by laser beam!", x, y+50, scale, scale, color.RGBA{0xff, 0xff, 0xff, 0xff})
	}

	return false
}

func newTowerIcon(x, y int) *icon {
	img, _, err := image.Decode(bytes.NewReader(towerImageData))
	if err != nil {
		log.Fatal(err)
	}

	return newIcon(x, y, ebiten.NewImageFromImage(img))
}

func (b *tower) Health() int {
	return b.health
}

func (b *tower) IsClicked(x, y int) bool {
	w, h := b.Size()
	return b.x-w/2 <= x && x <= b.x+w/2 && b.y-h/2 <= y && y <= b.y+h/2
}

func (b *tower) SetOverlap(overlap bool) {
	b.isOverlapping = overlap
}

func (b *tower) IsOverlap() bool {
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

func (b *tower) Cost() int {
	return CostTowerBuild
}
