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
	running bool
)

func init() {
	_, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(pling)))
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func playmp3() {
	if running {
		return
	}

	// Only run between 6:00 and 20:00
	now := time.Now()
	if now.Hour() < 6 || now.Hour() > 20 {
		return
	}

	streamer, _, err := mp3.Decode(io.NopCloser(bytes.NewReader(pling)))
	if err != nil {
		log.Fatal(err)
	}

	running = true

	log.Printf("User afk\n")

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		// This callback runs after streamer completed playing mp3
		running = false
	})))
}

func main() {
	idle := xidle.Idlemon{
		IdleOver: playmp3,
		// Play sound when user was idle for over 10 mins
		IdleOverT: time.Minute * 10,
		// Play sound every minute
		PollT: time.Minute,
	}
	// testing
	/*
	  idle.IdleOverT = time.Second
	  idle.PollT = time.Second * 5
	*/

	idle.Run()
}
