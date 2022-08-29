package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

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

func flagOverride(cCtx *cli.Context) {
	if cCtx.Duration("afk_timeout") != 0 {
		config.viper.Set("afk_timeout", cCtx.Duration(SchemaFieldToEnvName("AfkTimeout")))
	}

	if cCtx.IsSet("i3lock") {
		config.viper.Set("i3lock", cCtx.Bool(SchemaFieldToEnvName("I3lock")))
	}

	if cCtx.String("i3lock_color") != "" {
		config.viper.Set("i3lock_color", cCtx.String(SchemaFieldToEnvName("I3lockColor")))
	}

	if cCtx.String("mp3_file") != "" {
		config.viper.Set("mp3_file", cCtx.String(SchemaFieldToEnvName("Mp3File")))
	}

	if cCtx.IsSet("mp3_hour_start") {
		config.viper.Set("mp3_hour_start", cCtx.Int(SchemaFieldToEnvName("Mp3HourStart")))
	}

	if cCtx.IsSet("mp3_hour_stop") {
		config.viper.Set("mp3_hour_stop", cCtx.Int(SchemaFieldToEnvName("Mp3HourStop")))
	}

	if cCtx.Duration("mp3_interval") != 0 {
		config.viper.Set("mp3_interval", cCtx.Duration(SchemaFieldToEnvName("Mp3Interval")))
	}
}

func run(cCtx *cli.Context) error {
	var err error
	// Load config
	config, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	flagOverride(cCtx)

	if cCtx.IsSet("info") {
		log.Println(config)
		os.Exit(0)
	}

	mp3player = &Mp3Player{
		File:      config.Mp3File(),
		HourStart: config.Mp3HourStart(),
		HourStop:  config.Mp3HourStop(),
	}
	mp3player.Init()

	lastSeen = time.Now()
	lastPresent = time.Now()

	idle := xidle.Idlemon{
		IdleOver:        afk,
		IdleLess:        present,
		IdleOverTimeout: config.AfkTimeout(),
		IdleLessTimeout: time.Second,
	}

	idle.Run()
	return nil
}
func main() {
	app := &cli.App{
		Name:    "Pling",
		Usage:   "CLI for Pling",
		Version: "v0.1.0",

		Flags: []cli.Flag{
			&cli.DurationFlag{Name: "afk_timeout", Usage: "Duration until user is afk"},
			&cli.BoolFlag{Name: "i3lock", Usage: "Exec i3lock when afk_timeout is reached"},
			&cli.StringFlag{Name: "i3lock_color", Usage: "Change i3lock screen color to rrggbb"},
			&cli.StringFlag{Name: "mp3_file", Usage: "mp3 file played when user is afk"},
			&cli.IntFlag{Name: "mp3_hour_start", Usage: "Hour of day to start playing mp3"},
			&cli.IntFlag{Name: "mp3_hour_stop", Usage: "Hour of day to stop playing mp3"},
			&cli.DurationFlag{Name: "mp3_interval", Usage: "Duration to wait until playing mp3 file again"},
			&cli.BoolFlag{Name: "info", Usage: "Show config and exit"},
		},
		Action: run,
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
