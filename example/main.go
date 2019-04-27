package main

import (
	"github.com/wailovet/osmanthuswine"
	"github.com/wailovet/osmanthuswine/example/app"
	"github.com/wailovet/osmanthuswine/src/core"
)

func main() {
	core.GetInstanceRouterManage().Registered(&app.Wstest{})
	core.GetInstanceRouterManage().Registered(&app.Test{})
	osmanthuswine.Run()
}
