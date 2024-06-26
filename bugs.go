package main

import (
	"bytes"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	_ "embed"
	"image/color"
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

	// 攻撃中であるかどうかを示すフラグと攻撃アニメーションの時間
	// 攻撃アニメーションを行うために用いる
	attacking            bool
	attackDuration       int
	originalX, originalY int

	// 死亡時のアニメーションを管理するための変数
	deadAnimationDuration int

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

	bug := &bug{
		game: game,

		x:         x,
		y:         y,
		width:     bugImage.Bounds().Dx(),
		height:    bugImage.Bounds().Dy(),
		zindex:    50,
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

	switch bugColor {
	case bugsRed:
		bug.speed = 5
		bug.attackPower = 1
		bug.attackRange = 1
		bug.health = 2
		bug.name = "Red bug"
	case bugsBlue:
		bug.speed = 4
		bug.attackPower = 1
		bug.attackRange = 1
		bug.health = 3
		bug.name = "Blue bug"
	case bugsGreen:
		bug.speed = 3
		bug.attackPower = 1
		bug.attackRange = 50
		bug.health = 4
		bug.name = "Green bug"
	default:
		log.Fatal("invalid bug color")
	}

	return bug
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
	if b.health <= 0 {
		b.deadAnimationDuration++
		if b.deadAnimationDuration >= 30 {
			b.onDestroy(b)
		}
		return
	}

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

	// エフェクトを表示する
	switch b.selfColor {
	case bugsRed:
		// TODO: implement
	case bugsBlue:
		// TODO: implement
	case bugsGreen:
		tx, ty := a.(Building).Position()
		e := newGreenBugAttackEffect(b.game, b.x, b.y, tx, ty)
		b.game.updateHandler.Add(e)
		b.game.drawHandler.Add(e)
	}
}

type greenBugAttackEffect struct {
	game *Game

	currentX, currentY int
	startX, startY     int
	targetX, targetY   int

	erapsedFrame int

	zindex int
}

func newGreenBugAttackEffect(game *Game, startX, startY, targetX, targetY int) *greenBugAttackEffect {
	return &greenBugAttackEffect{
		game: game,

		currentX: startX,
		currentY: startY,
		startX:   startX,
		startY:   startY,
		targetX:  targetX,
		targetY:  targetY,

		zindex: 220,
	}
}

func (e *greenBugAttackEffect) Update() {
	// 攻撃エフェクトを描画する
	// 直線上に移動するエフェクトを描画する
	// 30 frame かけて target に向かって移動する
	e.currentX = e.startX + (e.targetX-e.startX)*e.erapsedFrame/30
	e.currentY = e.startY + (e.targetY-e.startY)*e.erapsedFrame/30
	if e.erapsedFrame <= 30 {
		e.erapsedFrame++
		return
	}
	e.game.updateHandler.Remove(e)
	e.game.drawHandler.Remove(e)
}

func (e *greenBugAttackEffect) Draw(screen *ebiten.Image) {
	// 直線上に移動するエフェクトを描画する
	vector.StrokeLine(screen, float32(e.startX), float32(e.startY), float32(e.currentX), float32(e.currentY), 10, color.RGBA{R: 128, G: 128, B: 128, A: 128}, true)
}

func (e *greenBugAttackEffect) ZIndex() int {
	return e.zindex
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

			b.attacking = true
			b.attackDuration = 7
			b.originalX, b.originalY = b.x, b.y
		} else {
			// クールダウンを消化する
			b.attackCooldown -= 1
		}

		// 攻撃中であればアニメーション動作を行う
		if b.attacking {
			b.attackDuration--
			if b.attackDuration <= 0 {
				b.x = b.originalX
				b.y = b.originalY
				b.attacking = false
			} else {
				// 攻撃対象に向かって一瞬スプライトを移動させる
				targetX, targetY := attackTarget.(Building).Position()
				dx := (targetX - b.x) / 4
				dy := (targetY - b.y) / 4
				b.x += dx / b.attackDuration
				b.y += dy / b.attackDuration
			}
		}

		return
	}

	// house に向かう
	var moveTargetX, moveTargetY int
	var found bool

	for _, building := range b.game.buildings {
		if building.Name() == "House" {
			moveTargetX, moveTargetY = building.Position()
			found = true
			break
		}
	}

	if !found {
		// すべての建物が破壊されている場合はその場にとどまる
		return
	}

	// ターゲットへの直線距離を計算
	dx := moveTargetX - b.x
	dy := moveTargetY - b.y

	// 移動方向のラジアンを計算
	angle := math.Atan2(float64(dy), float64(dx))

	// 回避動作
	// 虫同士がぴったり重ならないようにするための計算
	// やや自信のないロジックではある
	avoidX, avoidY := 0.0, 0.0
	for _, e := range b.game.enemies {
		ee := e.(*bug)
		if ee != b {
			distX := float64(ee.x - b.x)
			distY := float64(ee.y - b.y)
			distance := math.Sqrt(distX*distX + distY*distY)
			if distance > 0 && distance < float64(b.width) {
				avoidX -= distX / distance
				avoidY -= distY / distance
			}
		}
	}

	// 移動
	moveX := math.Cos(angle)*b.speed + avoidX
	moveY := math.Sin(angle)*b.speed + avoidY
	b.x += int(moveX)
	b.y += int(moveY)
}

func blueBugUpdate(b *bug) {
	// 青虫の特徴
	// 最寄りの障害物に向かって進む。障害物にぶつかったら、ぶつかったものに対して攻撃を行う。
	// 攻撃は一定時間ごとに行う。攻撃機範囲はせまい。自身の周囲ちょっとくらい (赤虫と同じ)。
	// 体力は赤虫よりもちょっと多い。
	// 赤虫より多く出現する。
	// 動きの速さは普通。

	// 最寄りの障害物を探す
	var nearestBuilding Damager
	nearestDistance := math.MaxFloat64
	for _, building := range b.game.buildings {
		x, y := building.Position()
		// 対象の建物と bug の距離を計算
		dx := x - b.x
		dy := y - b.y
		distance := math.Sqrt(float64(dx*dx + dy*dy))
		if distance < nearestDistance {
			nearestDistance = distance
			nearestBuilding = building.(Damager)
		}
	}

	// 最寄りの建物が攻撃範囲内にあるか確認

	// 対象の建物と bugs の攻撃範囲を踏まえた当たり判定を行う
	// bugs は size + attackRange の範囲を当たり判定として用いる
	if nearestBuilding == nil {
		// すべての建物が破壊されている場合はその場にとどまる
		return
	}
	target := nearestBuilding.(Building)
	x, y := target.Position()
	width, height := target.Size()
	var attackTarget Damager
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
		attackTarget = nearestBuilding
	}

	if attackTarget != nil {
		// クールダウン中でなければ攻撃
		if b.attackCooldown <= 0 {
			b.attack(nearestBuilding)
			b.attackCooldown = 60

			b.attacking = true
			b.attackDuration = 7
			b.originalX, b.originalY = b.x, b.y
		} else {
			// クールダウンを消化する
			b.attackCooldown -= 1
		}

		// 攻撃中であればアニメーション動作を行う
		if b.attacking {
			b.attackDuration--
			if b.attackDuration <= 0 {
				b.x = b.originalX
				b.y = b.originalY
				b.attacking = false
			} else {
				// 攻撃対象に向かって一瞬スプライトを移動させる
				targetX, targetY := attackTarget.(Building).Position()
				dx := (targetX - b.x) / 4
				dy := (targetY - b.y) / 4
				b.x += dx / b.attackDuration
				b.y += dy / b.attackDuration
			}
		}

		// クールダウン中でかつ攻撃対象が攻撃範囲内にいるときにはその場にとどまる
		return
	}

	// 最寄りの建物に向かって移動
	moveTargetX, moveTargetY := nearestBuilding.(Building).Position()

	// ターゲットへの直線距離を計算
	dx := moveTargetX - b.x
	dy := moveTargetY - b.y

	// 移動方向のラジアンを計算
	angle := math.Atan2(float64(dy), float64(dx))

	// 回避動作
	// 虫同士がぴったり重ならないようにするための計算
	// やや自信のないロジックではある
	avoidX, avoidY := 0.0, 0.0
	for _, e := range b.game.enemies {
		ee := e.(*bug)
		if ee != b {
			distX := float64(ee.x - b.x)
			distY := float64(ee.y - b.y)
			distance := math.Sqrt(distX*distX + distY*distY)
			if distance > 0 && distance < float64(b.width) {
				avoidX -= distX / distance
				avoidY -= distY / distance
			}
		}
	}

	// 移動
	moveX := math.Cos(angle)*b.speed + avoidX
	moveY := math.Sin(angle)*b.speed + avoidY
	b.x += int(moveX)
	b.y += int(moveY)
}

func greenBugUpdate(b *bug) {
	// 緑虫の特徴
	// 家に向かって一直線に進む。
	// 攻撃は一定時間ごとに行う。攻撃範囲が広い。飛び道具のようなものを放つ。
	// 攻撃範囲に任意の障害物が入ったとき、その障害物に向かって攻撃を行う。
	// 体力は青虫よりも多い。
	// 出現頻度は低い。
	// 動きは遅い。

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
	var found bool
	for _, building := range b.game.buildings {
		if building.Name() == "House" {
			moveTargetX, moveTargetY = building.Position()
			found = true
			break
		}
	}

	if !found {
		// すべての建物が破壊されている場合はその場にとどまる
		return
	}

	// ターゲットへの直線距離を計算
	dx := moveTargetX - b.x
	dy := moveTargetY - b.y

	// 移動方向のラジアンを計算
	angle := math.Atan2(float64(dy), float64(dx))

	// 回避動作
	// 虫同士がぴったり重ならないようにするための計算
	// やや自信のないロジックではある
	avoidX, avoidY := 0.0, 0.0
	for _, e := range b.game.enemies {
		ee := e.(*bug)
		if ee != b {
			distX := float64(ee.x - b.x)
			distY := float64(ee.y - b.y)
			distance := math.Sqrt(distX*distX + distY*distY)
			if distance > 0 && distance < float64(b.width) {
				avoidX -= distX / distance
				avoidY -= distY / distance
			}
		}
	}

	// 移動
	moveX := math.Cos(angle)*b.speed + avoidX
	moveY := math.Sin(angle)*b.speed + avoidY
	b.x += int(moveX)
	b.y += int(moveY)
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
	if b.health <= 0 {
		// 死亡時のアニメーションを行う
		// ぺちゃんこになるように縮小する
		scale := 1.0 - float64(b.deadAnimationDuration)/15
		if scale < 0 {
			scale = 0
		}

		opts.GeoM.Translate(0, float64(-b.height)/2)
		opts.GeoM.Scale(1, scale)
		opts.GeoM.Translate(0, float64(b.height)/2)
	} else {
		opts.GeoM.Scale(b.scale, b.scale)
	}
	opts.GeoM.Translate(float64(b.x)-float64(b.width)*b.scale/2, float64(b.y)-float64(b.height)*b.scale/2)
	screen.DrawImage(b.image, opts)
}

func (b *bug) ZIndex() int {
	return b.zindex
}

func (b *bug) OnClick(x, y int) bool {
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

	return false
}

func (b *bug) Health() int {
	return b.health
}

func (b *bug) IsClicked(x, y int) bool {
	width, height := b.width, b.height
	return b.x-width/2 <= x && x <= b.x+width/2 && b.y-height/2 <= y && y <= b.y+height/2
}
