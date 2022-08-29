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

type Mp3Player struct {
	File      string
	HourStart int
	HourStop  int

	playing bool

	// finishedPlayingTimestamp is updated after playing mp3
	// durationSinceLastPlay describes duration since last time finishing playing mp3
	// Check if this duration is greater than the config config.Mp3Interval()
	finishedPlayingTimestamp time.Time
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

// Plays mp3 when certain conditions are met
// The callback for speaker.Play is asynchronous so always set m.playing=false when returning early
func (m *Mp3Player) Play() {
	now := time.Now()

	if m.playing {
		return
	}

	// Compare duration since last mp3 ended until now and compare it to the config.Mp3Interval()
	if now.Sub(m.finishedPlayingTimestamp) < config.Mp3Interval() {
		m.playing = false
		return
	}

	// config.Mp3HourStart() and config.Mp3HourStop() allow configuring the hours of day between which to play mp3
	if now.Hour() < m.HourStart {
		if !notifiedAfk {
			log.Printf("playmp3(): Skip: only playing mp3 from %d o'clock\n", m.HourStart)
		}
		m.playing = false
		return
	}

	if now.Hour() != 0 && now.Hour() >= m.HourStop {
		if !notifiedAfk {
			log.Printf("playmp3(): Skip: only playing mp3 until %d o'clock\n", m.HourStop)
		}
		m.playing = false
		return
	}

	// All good - play mp3
	m.playing = true
	streamer, _, err := mp3.Decode(m.GetMp3ReadCloser())
	if err != nil {
		log.Fatal(err)
	}
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		// This callback runs after streamer completed playing mp3
		m.playing = false
		m.finishedPlayingTimestamp = time.Now()
	})))
}
