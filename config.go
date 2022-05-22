package main

import "github.com/spf13/viper"

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("tubelas")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/tubelas/")
	viper.AddConfigPath("$HOME/.config/tubelas/")
	viper.AddConfigPath(".")

	viper.SetDefault("http.listen", ":8080")
	viper.SetDefault("db.uri", "dbname=tubelas")
	viper.SetDefault("db.max_conns", 0)
	viper.SetDefault("log.format", "cli")
}

func loadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}
