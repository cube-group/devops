package setting

import (
	"app/library/log"
)

var Version string

func initLocal() {
	//local/config.yaml
	result, err := EmbedLocal().ReadFile("local/version")
	if err == nil {
		Version = string(result)
	}
	if err != nil {
		log.StdFatal("init", "initLocal", err)
	}
}
