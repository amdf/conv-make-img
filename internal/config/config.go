package config

import (
	"github.com/BurntSushi/toml"
)

//FileName где лежит.
const FileName = "configs/config.toml"

var sc *Config

type Config struct {
	Conv Converter `toml:"converter"`
}

type Converter struct {
	Host string
	Port string
}

//Load from file
func Load() (err error) {
	sc = new(Config)
	_, err = toml.DecodeFile(FileName, sc)

	return err
}

func GetConverterAddress() (result string) {
	if nil == sc {
		return
	}
	cfg := sc.Conv
	return cfg.Host + ":" + cfg.Port
}
