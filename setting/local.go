package setting

import (
	"app/library/log"
)

var Version string
var UpgradeLog string

func initLocal() {
	//local/config.yaml
	if result, err := EmbedLocal().ReadFile("local/version"); err == nil {
		Version = string(result)
	} else {
		log.StdFatal("init", "initLocal version", err)
	}
	//local/upgrade.html
	if result, err := EmbedLocal().ReadFile("local/upgrade.html"); err == nil {
		UpgradeLog = string(result)
	} else {
		log.StdFatal("init", "initLocal upgrade.html", err)
	}
}
