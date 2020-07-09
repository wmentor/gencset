package main

import (
	"flag"
	"runtime/debug"

	"github.com/spf13/viper"
	"github.com/wmentor/daemon"
	"github.com/wmentor/gencset/db"
	"github.com/wmentor/log"
	"github.com/wmentor/serv"

	_ "github.com/wmentor/gencset/controller"
)

func main() {

	var runDaemon bool
	flag.BoolVar(&runDaemon, "d", false, "run as daemon")
	flag.Parse()

	viper.SetConfigName("gencset")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if runDaemon {
		daemon.Run(viper.GetString("daemon"))
	}

	log.Open(viper.GetString("log"))

	if gc := viper.GetInt("garbage"); gc > 0 && gc < 100 {
		log.Infof("gc=%d", gc)
		debug.SetGCPercent((gc))
	}

	db.SetConnectString(viper.GetString("database"))

	serv.LoadTemplates(viper.GetString("templates"))

	serv.SetLogger(func(ld *serv.LogData) {
		format := "%s %s %d \"%s\" \"%s\" %.3f"
		if ld.StatusCode >= 500 {
			log.Errorf(format, ld.Addr, ld.Method, ld.StatusCode, ld.RequestURL, ld.UserAgent, ld.Seconds)
		} else {
			log.Infof(format, ld.Addr, ld.Method, ld.StatusCode, ld.RequestURL, ld.UserAgent, ld.Seconds)
		}
	})

	serv.SetErrorHandler(func(err error) {
		log.Error(err.Error())
	})

	if err := serv.Start(viper.GetString("server")); err != nil {
		log.Fatal(err.Error())
	}
}
