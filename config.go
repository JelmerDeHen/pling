package main

import (
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type ConfigSchema struct {
	// When I3lock is true i3lock will be excuted when user is idle
	I3lock bool `env:"i3lock" default:false`
	// Background color for i3
	I3lockColor string `env:"i3lock_color" default:"000000"`
	// Time until afk() handler is triggered
	IdleOverTimeout time.Duration `env:"idle_over_timeout" default:"10m"`
	// Not implemented
	Mp3 string `env:"mp3" default:""`
	// Hour of day to start playing mp3
	Mp3HourStart int `env:"mp3_hour_start" default:0`
	// Hour of day to stop playing mp3
	Mp3HourStop int `env:"mp3_hour_stop" default:0`
	// Time to wait between playing MP3
	Mp3Interval time.Duration `env:"mp3_interval" default:"5s"`
}

func SchemaFieldToEnvName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, _ := schemaT.FieldByName(field)
	return item.Tag.Get("env")
}

type ConfigI interface {
	I3lock() bool
	I3lockColor() string
	IdleOverTimeout() time.Duration
	Mp3() string
	Mp3HourStart() int
	Mp3HourStop() int
	Mp3Interval() time.Duration
}

type Config struct {
	viper *viper.Viper
	ConfigI
}

func (c *Config) I3lock() bool {
	return c.viper.GetBool(SchemaFieldToEnvName("I3lock"))
}

func (c *Config) I3lockColor() string {
	color := c.viper.GetString(SchemaFieldToEnvName("I3lockColor"))
	// rrggbb
	for i, c := range color {
		if !((c > 0x2f && c < 0x3a) || (c > 0x60 && c < 0x67) || (c > 0x40 && c < 0x47)) {
			return "000000"
		}
		if i > 5 {
			return "000000"
		}
	}
	return color
}

func (c *Config) IdleOverTimeout() time.Duration {
	return c.viper.GetDuration(SchemaFieldToEnvName("IdleOverTimeout"))
}

func (c *Config) Mp3() string {
	return c.viper.GetString(SchemaFieldToEnvName("Mp3"))
}

func (c *Config) Mp3HourStart() int {
	start := c.viper.GetInt(SchemaFieldToEnvName("Mp3HourStart"))

	return start
}
func (c *Config) Mp3HourStop() int {
	stop := c.viper.GetInt(SchemaFieldToEnvName("Mp3HourStop"))
	return stop
}

func (c *Config) Mp3Interval() time.Duration {
	return c.viper.GetDuration(SchemaFieldToEnvName("Mp3Interval"))
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetEnvPrefix("pling")
	v.AddConfigPath("$HOME")
	v.AddConfigPath(".")
	v.SetConfigName(".pling")

	// Load ConfigSchema to viper
	schemaT := reflect.TypeOf(ConfigSchema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		name := item.Tag.Get("env")
		def, exists := item.Tag.Lookup("default")
		if exists {
			v.SetDefault(name, def)
		}
		v.BindEnv(name)
	}

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c := &Config{
		viper: v,
	}

	return c, nil
}
