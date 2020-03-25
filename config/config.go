package config

import (
	"github.com/BurntSushi/toml"
)

// PathConfig represents system path configuration information.
type PathConfig struct {
	OutputCompareDir string `toml:"output_compare_dir"`
}

// Config represents application configuration.
type Config struct {
	PATH PathConfig `toml:"path"`
}

const confDir = "/go/src/work/config/" //　設定ファイルへの実行ファイルからの相対パスを指定

// NewConfig return configuration struct.
func NewConfig(appMode string) (Config, error) {
	var conf Config

	confPath := confDir + appMode + ".toml"
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
