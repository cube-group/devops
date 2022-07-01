package setting

import (
	"app/library/log"
	"bytes"
	"github.com/spf13/viper"
)

var Local *viper.Viper

func initLocal() {
	//local/config.yaml
	result, err := EmbedLocal().ReadFile("local/config.yaml")
	if err == nil {
		vip := viper.New()
		vip.SetConfigType("yaml")
		if err = vip.ReadConfig(bytes.NewBuffer(result)); err == nil {
			Local = vip
		}
	}
	if err != nil {
		log.StdFatal("init", "initLocal", err)
	}
}
