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

	// Only run between 6:00 and 20:00
	/*
	     now := time.Now()
	   	if now.Hour() < 6 || now.Hour() > 20 {
	   		return
	   	}
	*/

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
	log.Printf("Pling: idleOverTimeout=%s; pollInterval=%s; mp3=%s\n", config.IdleOverTimeout(), config.PollInterval(), config.Mp3())

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
