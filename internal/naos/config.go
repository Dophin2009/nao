package naos

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"gitlab.com/Dophin2009/nao/internal/config"
)

// Configuration contains config properties read from config files.
type Configuration struct {
	Hostname string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
	DB       struct {
		Path     string `mapstructure:"path"`
		Filemode uint32 `mapstructure:"filemode"`
	} `mapstructure:"db"`
}

// ReadConfigs returns a Configuration object with configuration properties
// read from standard directories.
func ReadConfigs() (*Configuration, error) {
	filename := "naos"

	var conf Configuration
	err := config.ReadConfigs(filename, ConfigDirs(), &conf)
	if err != nil {
		return nil,
			fmt.Errorf("failed to read config files %q: %w", filename+".*", err)
	}

	return &conf, nil
}

// ConfigDirs returns a list of configuration directories.
func ConfigDirs() []string {
	subdir := "nao"
	dirs := append(xdg.ConfigDirs, xdg.ConfigHome)
	for i := range dirs {
		dirs[i] = filepath.Join(dirs[i], subdir)
	}
	return dirs
}
