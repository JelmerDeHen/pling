package main

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"github.com/JelmerDeHen/xidle"
)

var (
	config  *Config
	running bool

	notifiedAfk     bool
	notifiedPresent bool

	lastSeen    time.Time
	lastPresent time.Time

	afkTime     time.Duration
	presentTime time.Duration

	finishedPlayingTimestamp time.Time
	playMp3Interval          time.Duration
)

func init() {
	lastSeen = time.Now()
	lastPresent = time.Now()

	finishedPlayingTimestamp = time.Now()
	playMp3Interval = time.Second * 5

	_, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(pling)))
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	config, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func afk() {
	if running {
		return
	}
	running = true
	notifiedPresent = false

	playmp3()
	if !notifiedAfk {

		if config.I3lock() {
			i3lock(config.I3lockColor())
		}

		lastSeen = time.Now()
		presentTime = lastSeen.Sub(lastPresent)

		log.Printf("User afk: presentTime=%s\n", presentTime)
		notifiedAfk = true
	}

}

func present() {
	notifiedAfk = false
	if !notifiedPresent {

		lastPresent := time.Now()
		afkTime = lastPresent.Sub(lastSeen)

		log.Printf("User is back: afkTime=%s\n", afkTime)
		notifiedPresent = true
	}
}

// Plays mp3 when certain conditions are met
// The callback for speaker.Play is asynchronous so always set running=false when returning early
func playmp3() {
	now := time.Now()

	// finishedPlayingTimestamp is updated after playing mp3
	// durationSinceLastPlay describes duration since last time finishing playing mp3
	// Check if this duration is greater than the config config.Mp3Interval()
	durationSinceLastPlay := now.Sub(finishedPlayingTimestamp)
	if durationSinceLastPlay < config.Mp3Interval() {
		running = false
		return
	}

	// config.Mp3HourStart() and config.Mp3HourStop() allow configuring the time between which to play mp3
	start := config.Mp3HourStart()
	if now.Hour() < start {
		if !notifiedAfk {
			log.Printf("playmp3(): Skip: only playing mp3 from %d o'clock\n", start)
		}
		running = false
		return
	}

	stop := config.Mp3HourStop()
	if now.Hour() != 0 && now.Hour() > stop {
		if !notifiedAfk {
			log.Printf("playmp3(): Skip: only playing mp3 until %d o'clock\n", stop)
		}
		running = false
		return
	}

	streamer, _, err := mp3.Decode(io.NopCloser(bytes.NewReader(pling)))
	if err != nil {
		log.Fatal(err)
	}
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		// This callback runs after streamer completed playing mp3
		running = false
		finishedPlayingTimestamp = time.Now()
	})))
	return
}

func main() {
	log.Printf("Pling: idleOverTimeout=%s; mp3=%s; mp3Interval=%s; mp3HourStart=%d; mp3HourStop=%d\n", config.IdleOverTimeout(), config.Mp3(), config.Mp3Interval(), config.Mp3HourStart(), config.Mp3HourStop())

	idle := xidle.Idlemon{
		IdleOver: afk,
		// Determines afk duration until mp3 is played
		IdleOverTimeout: config.IdleOverTimeout(),
		IdleLessTimeout: time.Second,
		IdleLess:        present,
	}

	idle.Run()
}
