package audio

import (
	"log"

	"github.com/veandco/go-sdl2/sdl_mixer"
)

type AudioID int

const (
	SOUND_COIN = iota
	SOUND_KICK
	SOUND_FIREBALL
	SOUND_POWERUP
	SOUND_BREAK_BRICK
)

var sounds map[AudioID]*mix.Chunk = make(map[AudioID]*mix.Chunk)

func LoadAllAudios() {
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
}

func Destroy() {
	for _, chunk := range sounds {
		chunk.Free()
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PlaySound(id AudioID) {
	sounds[id].Play(-1, 0)
}
