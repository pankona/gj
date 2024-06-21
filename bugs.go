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

//go:embed assets/bugs.png
var bugsImageData []byte

type bugColor int

type bug struct {
	game *Game

	x, y          int
	width, height int
	zindex        int
	image         *ebiten.Image
	selfColor     bugColor

	name        string
	health      int
	speed       float64
	attackPower int
	attackRange float64

	// 攻撃クールダウン
	// 初期化時に設定するものではなく、攻撃後に設定するもの
	attackCooldown int

	// 画像の拡大率。
	// TODO: 本当は画像のサイズそのものを変更したほうが見た目も処理効率も良くなる。余裕があれば後々やろう。
	scale float64

	// health が 0 になったときに呼ばれる関数
	onDestroy func(b *bug)
}

const (
	bugsRed bugColor = iota
	bugsBlue
	bugsGreen
)

func newBug(game *Game, bugColor bugColor, x, y int, onDestroy func(b *bug)) *bug {
	img, _, err := image.Decode(bytes.NewReader(bugsImageData))
	if err != nil {
		log.Fatal(err)
	}

	bugsImage := ebiten.NewImageFromImage(img)
	rect := func() image.Rectangle {
		switch bugColor {
		case bugsRed:
			return redBug()
		case bugsBlue:
			return blueBug()
		case bugsGreen:
			return greenBug()
		}
		log.Fatal("invalid bug color")
		return image.Rectangle{}
	}()
	bugImage := bugsImage.SubImage(rect).(*ebiten.Image)

	return &bug{
		game: game,

		x:         x,
		y:         y,
		width:     bugImage.Bounds().Dx(),
		height:    bugImage.Bounds().Dy(),
		image:     bugImage,
		selfColor: bugColor,

		// TODO: 虫種別によって異なる値を設定できるようにする
		name:   "Red bug",
		health: 2,

		speed:       5,
		attackPower: 1,
		attackRange: 1,

		scale: 1,

		onDestroy: onDestroy,
	}
}

func redBug() image.Rectangle {
	return image.Rect(1, 5, 29, 45)
}

func blueBug() image.Rectangle {
	return image.Rect(36, 4, 65, 45)
}

func greenBug() image.Rectangle {
	return image.Rect(35, 50, 66, 96)
}

func (b *bug) Update() {
	switch b.selfColor {
	case bugsRed:
		redBugUpdate(b)
	case bugsBlue:
		blueBugUpdate(b)
	case bugsGreen:
		greenBugUpdate(b)
	default:
		log.Fatal("invalid bug color")
	}
}

func (b *bug) attack(a Damager) {
	a.Damage(b.attackPower)
}

type Damager interface {
	Damage(int)
}

func (b *bug) Damage(d int) {
	if b.health <= 0 {
		return
	}

	b.health -= d

	if b.health <= 0 {
		b.health = 0
		b.onDestroy(b)
	}
}

type rect struct {
	x, y, width, height int
}

func intersects(r1, r2 rect) bool {
	return r1.x < r2.x+r2.width &&
		r2.x < r1.x+r1.width &&
		r1.y < r2.y+r2.height &&
		r2.y < r1.y+r1.height
}

func redBugUpdate(b *bug) {
	// target に向かう途中に障害物が攻撃射程に入ったとき、その障害物を target とする
	// いずれかの建物が攻撃レンジに入っているか確認
	var attackTarget Damager
	for _, building := range b.game.buildings {
		x, y := building.Position()
		width, height := building.Size()

		// 対象の建物と bugs の攻撃範囲を踏まえた当たり判定を行う
		// bugs は size + attackRange の範囲を当たり判定として用いる
		if intersects(
			// bug
			rect{
				b.x - b.width/2 - int(b.attackRange), b.y - b.height/2 - int(b.attackRange),
				b.width + int(b.attackRange)*2, b.height + int(b.attackRange)*2,
			},
			// building
			rect{x - width/2, y - height/2,
				width, height},
		) {
			// 攻撃射程圏内であるので、その建物を attack 対象にする
			attackTarget = building.(Damager)

			break
		}
	}

	// attack target がいるならば攻撃する。そうでないならば house に向かう
	if attackTarget != nil {
		// クールダウン中でなければ攻撃
		if b.attackCooldown <= 0 {
			b.attack(attackTarget)
			b.attackCooldown = 60
		} else {
			// クールダウンを消化する
			b.attackCooldown -= 1
		}

		return
	}

	// house に向かう
	var moveTargetX, moveTargetY int

	for _, building := range b.game.buildings {
		if building.Name() == "House" {
			moveTargetX, moveTargetY = building.Position()
			break
		}
	}

	// ターゲットへの直線距離を計算
	dx := moveTargetX - b.x
	dy := moveTargetY - b.y

	// 移動方向のラジアンを計算
	angle := math.Atan2(float64(dy), float64(dx))

	// 移動
	b.x += int(math.Cos(angle) * b.speed)
	b.y += int(math.Sin(angle) * b.speed)
}

func blueBugUpdate(b *bug) {
	// todo: implement
}

func greenBugUpdate(b *bug) {
	// todo: implement
}

func (b *bug) Name() string {
	return b.name
}

func (b *bug) Position() (int, int) {
	return b.x, b.y
}

func (b *bug) Size() (int, int) {
	return b.width, b.height
}

// 画面中央に配置
func (b *bug) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(b.scale, b.scale)
	opts.GeoM.Translate(float64(b.x)-float64(b.width)*b.scale/2, float64(b.y)-float64(b.height)*b.scale/2)
	screen.DrawImage(b.image, opts)
}

func (b *bug) ZIndex() int {
	return b.zindex
}

func (b *bug) OnClick(x, y int) {
	switch b.selfColor {
	case bugsRed:
		b.game.clickedObject = "red bug"
	case bugsBlue:
		b.game.clickedObject = "blue bug"
	case bugsGreen:
		b.game.clickedObject = "green bug"
	default:
		log.Fatal("invalid bug color")
	}

	// infoPanel に情報を表示する
	icon := newBugIcon(80, eScreenHeight+70, b.selfColor)
	b.game.infoPanel.setIcon(icon)
	b.game.infoPanel.setUnit(b)
}

func (b *bug) Health() int {
	return b.health
}

func (b *bug) IsClicked(x, y int) bool {
	width, height := b.width, b.height
	return b.x-width/2 <= x && x <= b.x+width/2 && b.y-height/2 <= y && y <= b.y+height/2
}
