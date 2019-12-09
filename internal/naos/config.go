package naos

import (
	"github.com/OpenPeeDeeP/xdg"
	"github.com/spf13/viper"
)

// Configuration contains config properties
// read from config files.
type Configuration struct {
	Hostname string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
	DB       struct {
		Path     string `mapstructure:"path"`
		Filemode uint32 `mapstructure:"filemode"`
	} `mapstructure:"db"`
}

// ReadLinuxConfigs looks for config files from
// the default Linux config locations and returns
//a parsed Configuration struct
func ReadLinuxConfigs() (Configuration, error) {
	etcDir := "/etc/nao/"
	userDir := xdg.ConfigHome() + "/nao/"
	confFileDirs := []string{etcDir, userDir}
	return ReadConfig(confFileDirs)
}

// ReadConfig returns a Configuration struct
// parsed from config files in the given directories
func ReadConfig(filedirs []string) (conf Configuration, err error) {
	viper.SetConfigName("naos")
	for _, f := range filedirs {
		viper.AddConfigPath(f)
	}
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&conf)
	return
}
