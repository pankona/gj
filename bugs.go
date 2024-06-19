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
	game          *Game
	x, y          int
	width, height int
	zindex        int
	image         *ebiten.Image
	selfColor     bugColor

	speed       float64
	attackPower int
	attackRange float64

	// 攻撃クールダウン
	// 初期化時に設定するものではなく、攻撃後に設定するもの
	attackCooldown int

	// 画像の拡大率。
	// TODO: 本当は画像のサイズそのものを変更したほうが見た目も処理効率も良くなる。余裕があれば後々やろう。
	scale float64
}

const (
	bugsRed bugColor = iota
	bugsBlue
	bugsGreen
)

func newBug(game *Game, bugColor bugColor, x, y int) *bug {
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
		game:      game,
		x:         x,
		y:         y,
		width:     bugImage.Bounds().Dx(),
		height:    bugImage.Bounds().Dy(),
		image:     bugImage,
		selfColor: bugColor,

		speed:       5,
		attackPower: 1,
		attackRange: 5,

		scale: 1,
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

func redBugUpdate(b *bug) {
	// builds から house を探して target とする
	var (
		targetX, targetY          int
		targetWidth, targetHeight int
		target                    Damager
	)

	for _, building := range b.game.buildings {
		if building.Name() == "house" {
			targetX, targetY = building.Position()
			targetWidth, targetHeight = building.Size()
			target = building.(*house)
			break
		}
	}

	// 自分の中心座標を基準に行き先を計算する
	bx, by := b.x+b.width/2, b.y+b.height/2

	// ターゲットへの直線距離を計算
	dx := targetX - bx
	dy := targetY - by
	distance := math.Sqrt(float64(dx*dx + dy*dy))

	// 攻撃レンジに入っているか確認
	targetSize := targetWidth
	if targetWidth > targetHeight {
		targetSize = targetHeight
	}

	if distance <= b.attackRange+float64(targetSize)/2 {
		// クールダウン中でなければ攻撃
		if b.attackCooldown <= 0 {
			b.attack(target)
			b.attackCooldown = 60
		} else {
			// クールダウンを消化する
			b.attackCooldown -= 1
		}

		return
	}

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

// 画面中央に配置
func (b *bug) Draw(screen *ebiten.Image) {
	// 画像を描画
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(b.scale, b.scale)
	opts.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.image, opts)
}

func (b *bug) ZIndex() int {
	return b.zindex
}
