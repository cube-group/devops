package main

import (
	"app/models"
	"app/setting"
	"app/web"
)

func main() {
	setting.Init()
	models.Init()
	web.Init()
}
