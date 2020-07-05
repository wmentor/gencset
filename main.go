package main

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/viper"
	"github.com/wmentor/gencset/db"
	"github.com/wmentor/log"
	"github.com/wmentor/serv"
)

func main() {

	viper.SetConfigName("gencset")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.Open(viper.GetString("log"))

	if gc := viper.GetInt("garbage"); gc > 0 && gc < 100 {
		log.Info(fmt.Sprintf("gc=%d", gc))
		debug.SetGCPercent((gc))
	}

	db.SetConnectString(viper.GetString("database"))

	if err := serv.Start(viper.GetString("serv")); err != nil {
		log.Fatal(err.Error())
	}
}
