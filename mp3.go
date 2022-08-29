package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

/*
running is neccessary because speaker should finish playing mp3 until next mp3 is played
*/

type Mp3Player struct {
	File    string
	Running bool
}

func (m *Mp3Player) Init() {
	// speaker.Init() should be called once
	_, format, err := mp3.Decode(m.GetMp3ReadCloser())
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func (m *Mp3Player) GetMp3ReadCloser() io.ReadCloser {
	if m.File == "" {
		return io.NopCloser(bytes.NewReader(pling))
	}
	rc, err := os.Open(m.File)
	if err != nil {
		log.Fatal(err)
	}
	return rc
}

func (m *Mp3Player) OkToPlay() {
	now := time.Now()

	if m.Running {
		return
	}

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
}

func (m *Mp3Player) Play() {
	// Play mp3 file
	streamer, _, err := mp3.Decode(m.GetMp3ReadCloser())
	if err != nil {
		log.Fatal(err)
	}
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		// This callback runs after streamer completed playing mp3
		running = false
		finishedPlayingTimestamp = time.Now()
	})))
}
