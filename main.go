package main

import (
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("gencset")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
