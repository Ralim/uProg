package main

import (
	"github.com/ralim/uprog/config"
	"github.com/ralim/uprog/ui"
)

func main() {
	conf := config.Config{}
	conf.ParseFlags()
	uiConf := ui.NewUI(&conf)
	uiConf.RunUI()
}
