package main

import (
	"app/db"
	"app/models"
	"app/setting"
	"app/web"
	"flag"
)

func init() {
	setting.Init()
	models.Init()
}

func main() {
	var cmd string
	flag.StringVar(&cmd, "cmd", "web", "web|db")
	flag.Parse()
	switch cmd {
	case "db":
		db.Init()
	default:
		web.Init()
	}
}
