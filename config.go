package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
)

func initConfig() {
	// defaults
	viper.SetDefault("AdminLocation", "/admin/") // TODO : some helper to convert strings to paths ? ie: path("admin") -> "/admin/"
	viper.SetDefault("contentdir", "content")

	// config name and location
	viper.SetConfigName("config")

	viper.AddConfigPath("/etc/topiary")
	viper.AddConfigPath(".")

	// read config
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("No topiary config found. Using defaults.")
	}

	// watch config ; TODO : config to turn this on/off
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

func getConfig(s string) string {
	return viper.GetString(s)
}
