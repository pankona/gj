package main

import (
	"bytes"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"

	_ "embed"
)

type audioPlayer struct {
	audioContext *audio.Context

	players map[string]*audio.Player
}

var aplayer *audioPlayer

var (
	//go:embed assets/bakuhatsu.mp3
	soundBakuhatsuMP3 []byte
	//go:embed assets/beam.mp3
	soundBeamMP3 []byte
	//go:embed assets/binta.mp3
	soundBintaMP3 []byte
	//go:embed assets/choice.mp3
	soundChoiceMP3 []byte
	//go:embed assets/clear.mp3
	soundClearMP3 []byte
	//go:embed assets/don.mp3
	soundDonMP3 []byte
	//go:embed assets/gameover.mp3
	soundGameoverMP3 []byte
	//go:embed assets/gyaa.mp3
	soundGyaaMP3 []byte
	//go:embed assets/hikkaki.mp3
	soundHikkakiMP3 []byte
	//go:embed assets/kettei.mp3
	soundKetteiMP3 []byte
	//go:embed assets/kuzureru.mp3
	soundKuzureruMP3 []byte
	//go:embed assets/shot.mp3
	soundShotMP3 []byte

	//go:embed assets/bgm.mp3
	soundBgmMP3 []byte
)

const (
	soundGyaa      = "gyaa"      // 虫が死んだときの音
	soundBakuhatsu = "bakuhatsu" // 電波塔が範囲攻撃するときの音
	soundBeam      = "beam"      // 塔がビームを撃つときの音
	soundBinta     = "binta"     // 手で叩くときの音
	soundChoice    = "choice"    // ボタンを押下したときの音
	soundClear     = "clear"     // ゲームクリア時の音
	soundDon       = "don"       // 建物を置いたときの音
	soundGameover  = "gameover"  // ゲームオーバー時の音
	soundHikkaki   = "hikkaki"   // 虫が建物を攻撃するときの音
	soundKettei    = "kettei"    // Ready ボタンを押したときの音
	soundKuzureru  = "kuzureru"  // 建物が壊れるときの音
	soundShot      = "shot"      // 緑虫が弾を撃つときの音

	soundBgm = "bgm"
)

func getAudioPlayer() *audioPlayer {
	if aplayer != nil {
		return aplayer
	}

	aplayer = &audioPlayer{
		audioContext: audio.NewContext(44100),
		players:      map[string]*audio.Player{},
	}

	sounds := map[string][]byte{
		soundGyaa:      soundGyaaMP3,
		soundBakuhatsu: soundBakuhatsuMP3,
		soundBeam:      soundBeamMP3,
		soundBinta:     soundBintaMP3,
		soundChoice:    soundChoiceMP3,
		soundClear:     soundClearMP3,
		soundDon:       soundDonMP3,
		soundGameover:  soundGameoverMP3,
		soundHikkaki:   soundHikkakiMP3,
		soundKettei:    soundKetteiMP3,
		soundKuzureru:  soundKuzureruMP3,
		soundShot:      soundShotMP3,
	}

	for name, buf := range sounds {
		stream, err := mp3.DecodeWithSampleRate(aplayer.audioContext.SampleRate(), bytes.NewReader(buf))
		if err != nil {
			log.Fatalf("failed to decode mp3: %v", err)
		}
		aplayer.players[name] = mustNewPlayer(aplayer.audioContext, stream)
	}

	stream, err := mp3.DecodeWithSampleRate(aplayer.audioContext.SampleRate(), bytes.NewReader(soundBgmMP3))
	if err != nil {
		log.Fatalf("failed to decode mp3: %v", err)
	}
	bgmStream := audio.NewInfiniteLoop(stream, 160_000*40.6)
	bgmPlayer := mustNewPlayer(aplayer.audioContext, bgmStream)
	aplayer.players[soundBgm] = bgmPlayer

	return aplayer
}

func mustNewPlayer(context *audio.Context, stream io.Reader) *audio.Player {
	player, err := context.NewPlayer(stream)
	if err != nil {
		log.Fatalf("failed to create player: %v", err)
	}
	player.SetVolume(0.2) // なんか音がうるさいのでちっさくしておく

	return player
}

func (a *audioPlayer) play(soundname string) {
	// 同時に鳴らす音を制限する
	var active int
	for _, player := range a.players {
		if player.IsPlaying() {
			active++
		}
	}
	if active >= 5 {
		return
	}

	if player, ok := a.players[soundname]; ok {
		if player.IsPlaying() {
			return
		}
		player.Rewind()
		player.Play()
		return
	}
}

func (a *audioPlayer) playBGM() {
	player := a.players[soundBgm]
	player.Rewind()
	player.Play()
}

func (a *audioPlayer) stopBGM() {
	player := a.players[soundBgm]
	player.Pause()
}
