package config

import (
	"github.com/BurntSushi/toml"
)

//FileName где лежит.
const FileName = "configs/config.toml"

var sc *Config

type Config struct {
	Converter ConverterCfg `toml:"converter"`
	Consumer  ConsumerCfg  `toml:"consumer"`
}

type ConverterCfg struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type ConsumerCfg struct {
	Topic   string   `toml:"topic"`
	Group   string   `toml:"group"`
	Brokers []string `toml:"brokers"`
	Verbose bool     `toml:"verbose"`
}

func Get() *Config {
	return sc
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
	cfg := sc.Converter
	return cfg.Host + ":" + cfg.Port
}
