package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	HomeAltitudeKM     float64 `required:"true" split_words:"true"`
	HomeLatitudeDeg    float64 `required:"true" split_words:"true"`
	HomeLongitudeDeg   float64 `required:"true" split_words:"true"`
	SpacetrackUsername string  `required:"true" split_words:"true"`
	SpacetrackPassword string  `required:"true" split_words:"true"`
}

func Get() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
