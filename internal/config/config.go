package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// ReadConfigs reads config files in the given directories
// with the given filename (without extension). The overall
// config is unmarshalled into the given pointer.
func ReadConfigs(filename string, dirs []string, structure interface{}) error {
	if structure == nil {
		return fmt.Errorf("structure: %w", errors.New("is nil"))
	}

	viper.SetConfigName(filename)
	for _, d := range dirs {
		viper.AddConfigPath(d)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read in configs: %w", err)
	}

	err = viper.Unmarshal(structure)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
