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
)

func init() {
	lastSeen = time.Now()
	lastPresent = time.Now()

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

func playmp3() {
	// Don't play sounds at night
	now := time.Now()
	start := config.Mp3HourStart()
	stop := config.Mp3HourStop()
	if stop != 0 {
		if now.Hour() < start || now.Hour() > stop {
			if !notifiedAfk {
				log.Printf("playmp3(): Skip: only playing mp3 between %d and %d\n", start, stop)
			}
			running = false
			return
		}
	}

	streamer, _, err := mp3.Decode(io.NopCloser(bytes.NewReader(pling)))
	if err != nil {
		log.Fatal(err)
	}
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		// This callback runs after streamer completed playing mp3
		running = false
	})))
}

func main() {
	log.Printf("Pling: idleOverTimeout=%s; pollInterval=%s; mp3=%s; mp3HourStart=%d; mp3HourStop=%d\n", config.IdleOverTimeout(), config.PollInterval(), config.Mp3(), config.Mp3HourStart(), config.Mp3HourStop())

	idle := xidle.Idlemon{
		IdleOver: afk,
		// Determines afk duration until mp3 is played
		IdleOverTimeout: config.IdleOverTimeout(),
		// Will determine the interval between mp3 plays
		PollInterval:    config.PollInterval(),
		IdleLessTimeout: time.Second * 1,
		IdleLess:        present,
	}

	idle.Run()
}
