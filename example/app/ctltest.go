package app

import (
	"log"

	"github.com/ouxuanserver/osmanthuswine/src/core"
)

type Ctltest struct {
	core.Controller
}

func (that *Ctltest) Prepare() {
	//该方法会在Method执行前调用，用户扩展类似校验登陆
	log.Println("Prepare")
}

func (that *Ctltest) Test() {
	that.DisplayByData("test")
}
