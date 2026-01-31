package internal

import (
	"bytes"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"

	"github.com/ihorbryk/manta/assets"
)

var (
	otoCtx  *oto.Context
	otoOnce sync.Once
)

// initOtoContext initializes the shared Oto context.
// This function should only be called once via sync.Once.
// Creating multiple contexts is NOT supported by the Oto library.
func initOtoContext() {
	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Create the context once and reuse it for all audio playback
	ctx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	otoCtx = ctx
}

func PlayNotification() {
	// Ensure the Oto context is initialized (only happens once)
	otoOnce.Do(initOtoContext)

	// Read the embedded mp3 file into memory
	fileBytes, err := assets.NotifySound.ReadFile("notify.mp3")
	if err != nil {
		panic("reading embedded notify.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Create a new 'player' that will handle our sound. Paused by default.
	// We reuse the shared context but create a new player for each playback.
	player := otoCtx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// Close the player to free resources after playback completes
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
