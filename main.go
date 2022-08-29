package main

import (
	"log"
	"time"

	"github.com/JelmerDeHen/xidle"
)

var (
	config *Config

	// These bools guard executing afk() and present() handler multiple times
	notifiedAfk     bool
	notifiedPresent bool

	lastSeen    time.Time
	lastPresent time.Time

	afkTime     time.Duration
	presentTime time.Duration

	finishedPlayingTimestamp time.Time
	playMp3Interval          time.Duration

	mp3player *Mp3Player
)

func init() {
	// Load config
	var err error
	config, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	mp3player = &Mp3Player{
		File:      config.Mp3File(),
		HourStart: config.Mp3HourStart(),
		HourStop:  config.Mp3HourStop(),
	}
	mp3player.Init()

	lastSeen = time.Now()
	lastPresent = time.Now()
}

func afk() {
	notifiedPresent = false
	mp3player.Play()

	if !notifiedAfk {
		if config.I3lock() {
			i3lock(config.I3lockColor())
		}

		now := time.Now()

		lastSeen = now
		presentTime = now.Sub(lastPresent)

		log.Printf("User afk: presentTime=%s\n", presentTime)
		notifiedAfk = true
	}
}

func present() {
	notifiedAfk = false
	if !notifiedPresent {
		now := time.Now()

		lastPresent = now
		afkTime = now.Sub(lastSeen)

		log.Printf("User is back: afkTime=%s\n", afkTime)
		notifiedPresent = true
	}
}

func main() {
	log.Println(config)

	idle := xidle.Idlemon{
		IdleOver: afk,
		// Determines afk duration until mp3 is played
		IdleOverTimeout: config.IdleOverTimeout(),
		IdleLessTimeout: time.Second,
		IdleLess:        present,
	}

	idle.Run()
}
