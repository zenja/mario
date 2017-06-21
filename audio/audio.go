package audio

import (
	"log"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl_mixer"
)

type SoundID int
type MusicID int

const (
	SOUND_COIN SoundID = iota
	SOUND_KICK
	SOUND_FIREBALL
	SOUND_POWERUP
	SOUND_BREAK_BRICK
	SOUND_HERO_DIE
	SOUND_PIPE
	SOUND_STOMP
	SOUND_BUMP
)

const (
	MUSIC_0 MusicID = iota
)

var sounds map[SoundID]*mix.Chunk = make(map[SoundID]*mix.Chunk)
var musics map[MusicID]*mix.Music = make(map[MusicID]*mix.Music)

var currentMusic MusicID = MUSIC_0

func InitAudio() {
	err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 2048)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to init audio system"))
	}

	loadAllAudios()
}

func loadAllAudios() {
	var err error
	sounds[SOUND_COIN], err = mix.LoadWAV("assets/audio/coin.wav")
	must(err)
	sounds[SOUND_KICK], err = mix.LoadWAV("assets/audio/kick.wav")
	must(err)
	sounds[SOUND_FIREBALL], err = mix.LoadWAV("assets/audio/fireball.wav")
	must(err)
	sounds[SOUND_POWERUP], err = mix.LoadWAV("assets/audio/powerup.wav")
	must(err)
	sounds[SOUND_BREAK_BRICK], err = mix.LoadWAV("assets/audio/break_brick.wav")
	must(err)
	sounds[SOUND_HERO_DIE], err = mix.LoadWAV("assets/audio/hero_die.wav")
	must(err)
	sounds[SOUND_PIPE], err = mix.LoadWAV("assets/audio/pipe.wav")
	must(err)
	sounds[SOUND_STOMP], err = mix.LoadWAV("assets/audio/stomp.wav")
	must(err)
	sounds[SOUND_BUMP], err = mix.LoadWAV("assets/audio/bump.wav")
	must(err)

	// music
	musics[MUSIC_0], err = mix.LoadMUS("assets/audio/music/mario-bg-music-0.wav")
	must(err)
}

func PlayMusic() {
	musics[currentMusic].Play(-1)
}

func PauseMusic() {
	mix.PauseMusic()
}

func StopMusic() {
	mix.HaltMusic()
}

func ReloadMusic(mid MusicID) {
	if mid == currentMusic {
		return
	}

	// stop current music
	mix.HaltMusic()
	musics[mid].Play(-1)
}

func Destroy() {
	for _, s := range sounds {
		s.Free()
	}
	for _, m := range musics {
		m.Free()
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PlaySound(id SoundID) {
	sounds[id].Play(-1, 0)
}
