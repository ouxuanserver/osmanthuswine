package main

import (
	"github.com/ouxuanserver/osmanthuswine"
	"github.com/ouxuanserver/osmanthuswine/example/app"
	"github.com/ouxuanserver/osmanthuswine/src/core"
)

func main() {
	core.GetInstanceRouterManage().Registered(&app.Wstest{})
	core.GetInstanceRouterManage().Registered(&app.Ctltest{})
	osmanthuswine.Run()
}
