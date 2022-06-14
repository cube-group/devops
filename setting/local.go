package setting

import (
	"app/library/log"
	"github.com/spf13/viper"
)

var Local *viper.Viper

func initLocal() {
	//local/config.yaml
	vip := viper.New()
	f, err := EmbedLocal().Open("local/config.yaml")
	if err != nil {
		goto ERROR
	}
	defer f.Close()
	err = vip.ReadConfig(f)
	if err != nil {
		goto ERROR
	}
	Local = vip
	return

ERROR:
	log.StdFatal("init", "initLocal", err)
}
