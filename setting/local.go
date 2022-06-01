package setting

import (
	"app/library/core"
	"github.com/spf13/viper"
	"log"
)

var Local *viper.Viper

func initLocal() {
	vip := viper.New()
	vip.SetConfigName("config")
	vip.SetConfigType("yaml")
	vip.AddConfigPath(".")
	if err := vip.ReadInConfig(); err != nil {
		log.Fatal("init", "initLocal", err)
	}
	core.Lock(func() {
		Local = vip
	})
}
