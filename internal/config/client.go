package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitClientConfig(cfgFile string) error {
	v := viper.New()
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)

		if err := v.ReadInConfig(); err != nil {
			cobra.CheckErr(fmt.Errorf("read config: %w", err))
		}
	}
	return nil
}
