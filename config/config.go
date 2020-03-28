package config

import (
	"github.com/BurntSushi/toml"
)

// PathConfig represents system path configuration information.
type PathConfig struct {
	OutputCompareDir string `toml:"output_compare_dir"`
}

// URLListConfig represents system path configuration information.
type URLListConfig struct {
	SourceURL string `toml:"source_url"`
	TargetURL string `toml:"target_url"`
}

// Config represents application configuration.
type Config struct {
	PATH    PathConfig      `toml:"path"`
	URLLIST []URLListConfig `toml:"items"`
}

const confDir = "/go/src/work/config/"

// NewConfig return configuration struct.
func NewConfig(appMode string) (Config, error) {
	var conf Config

	confPath := confDir + appMode + ".toml"
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
