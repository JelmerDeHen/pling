package main

import (
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type ConfigSchema struct {
	IdleOverTimeout time.Duration `env:"idle_over_timeout" default:"10m"`
	// Not implemented yet
	Mp3          string        `env:"mp3" default:""`
	PollInterval time.Duration `env:"poll_interval" default:"1s"`
}

func SchemaFieldToEnvName(field string) string {
	schemaT := reflect.TypeOf(ConfigSchema{})
	item, _ := schemaT.FieldByName(field)
	return item.Tag.Get("env")
}

type ConfigI interface {
	IdleOverTimeout() time.Duration
	PollInterval() time.Duration
	Mp3() string
}

type Config struct {
	viper *viper.Viper
	ConfigI
}

func (c *Config) IdleOverTimeout() time.Duration {
	return c.viper.GetDuration(SchemaFieldToEnvName("IdleOverTimeout"))
}

func (c *Config) PollInterval() time.Duration {
	return c.viper.GetDuration(SchemaFieldToEnvName("PollInterval"))
}

func (c *Config) Mp3() string {
	return c.viper.GetString(SchemaFieldToEnvName("Mp3"))
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
