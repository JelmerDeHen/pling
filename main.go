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
	config          *Config
	running         bool
	notifiedAfk     bool
	notifiedPresent bool
)

func init() {
	_, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(pling)))
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	config, err = NewConfig()
	if err != nil {
		log.Panic(err)
	}
}

func afk() {
	if running {
		return
	}
	running = true

	if !notifiedAfk {
		log.Println("User went afk")
		notifiedAfk = true

		if config.I3lock() {
			i3lock(config.I3lockColor())
		}

	}
	notifiedPresent = false

	playmp3()
}

func playmp3() {
	// Don't play sounds at night
	now := time.Now()
	start := config.Mp3HourStart()
	stop := config.Mp3HourStop()
	if stop != 0 {
		if now.Hour() < start || now.Hour() > stop {
			log.Printf("playmp3(): Skip: only playing mp3 between %d and %d\n", start, stop)
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

func idleLess() {
	if !notifiedPresent {
		log.Println("User not afk")
		notifiedPresent = true
	}
	notifiedAfk = false
}

func main() {
	log.Printf("Pling: idleOverTimeout=%s; pollInterval=%s; mp3=%s; mp3HourStart=%d; mp3HourStop=%d\n", config.IdleOverTimeout(), config.PollInterval(), config.Mp3(), config.Mp3HourStart(), config.Mp3HourStop())

	idle := xidle.Idlemon{
		IdleOver: afk,
		// Determines afk duration until mp3 is played
		IdleOverT: config.IdleOverTimeout(),
		// Will determine the interval between mp3 plays
		PollT:     config.PollInterval(),
		IdleLessT: time.Second * 1,
		IdleLess:  idleLess,
	}

	idle.Run()
}
