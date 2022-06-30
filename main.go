package main

import (
	"app/db"
	"app/models"
	"app/setting"
	"app/web"
	"embed"
	_ "embed"
	"flag"
)

//go:embed local
var embedLocal embed.FS

func init() {
	setting.Init(
		map[string]embed.FS{
			"local": embedLocal,
		})
	models.Init()
}

func main() {
	var cmd string
	flag.StringVar(&cmd, "cmd", "", "web|db|task")
	flag.Parse()
	switch cmd {
	case "db":
		db.Init()
	case "task":
		models.InitCronjob()
	case "web":
		web.Init()
	default:
		go models.InitCronjob()
		web.Init()
	}
}
