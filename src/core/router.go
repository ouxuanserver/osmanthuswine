package core

import (
	"github.com/wailovet/osmanthuswine/src/interfaces"
	"log"
	"reflect"
	"strings"
	"unicode"
)

type RouterManage struct {
	RegisteredData map[string]reflect.Type
}

var instanceRouterManage *RouterManage

func GetInstanceRouterManage() *RouterManage {
	if instanceRouterManage == nil {
		instanceRouterManage = &RouterManage{} // not thread safe
		instanceRouterManage.RegisteredData = make(map[string]reflect.Type)
	}
	return instanceRouterManage
}

func (rm *RouterManage) Registered(i interface{}) {
	t := reflect.ValueOf(i)
	GetInstanceRouterManage().RegisteredData[t.Type().String()] = reflect.Indirect(t).Type()
}

func (rm *RouterManage) GetModuleName(name string) string {
	if name == "" {
		return "index"
	}
	return strings.ToLower(name)
}
func (rm *RouterManage) GetControllerName(name string) string {
	if name == "" {
		return "Index"
	}
	for i, v := range name {
		return string(unicode.ToUpper(v)) + name[i+1:]
	}
	return "Index"
}
func (rm *RouterManage) GetFunName(name string) string {
	if name == "" {
		return "Index"
	}
	for i, v := range name {
		return string(unicode.ToUpper(v)) + name[i+1:]
	}
	return "Index"
}

func (rm *RouterManage) RouterSend(urlPath string, request Request, response Response, crossDomain string) {
	tmp := strings.Split(urlPath, ".")
	if len(tmp) > 1 {
		urlPath = strings.Join(tmp[0:len(tmp)-1], ".")
	}

	sar := strings.Split(urlPath, "/")
	for len(sar) < 5 {
		sar = append(sar, "")
	}
	//过滤非 /Api开头的
	module := rm.GetModuleName(sar[2])
	controller := rm.GetControllerName(sar[3])
	fun := rm.GetFunName(sar[4])

	ctr := "*" + module + "." + controller

	_, ok := rm.RegisteredData[ctr]
	if !ok {
		panic("未注册该组件:" + ctr)
	}

	vc := reflect.New(rm.RegisteredData[ctr])

	wsinit := vc.MethodByName("WebSocketInit")
	if wsinit.IsValid() {
		response.IsWebSocket = true
		hand := vc.Interface().(interfaces.WebSocketInterface)
		ws := GetWebSocket(ctr+"-"+fun, hand)
		hand.SetFunName(fun)
		hand.WebSocketInit(ws)
		defer func() {
			errs := recover()
			if errs == nil {
				return
			}
			log.Printf("websocket error:%v", errs)
		}()

		_ = ws.HandleRequest(response.OriginResponseWriter, request.OriginRequest)
		return
	}

	f := vc.MethodByName(fun)
	if !f.IsValid() {
		panic("组件找不到相应function:" + fun)
	}

	init := vc.MethodByName("ControllerInit")
	if init.IsValid() {
		init.Call([]reflect.Value{reflect.ValueOf(request), reflect.ValueOf(response)})
		f.Call(nil)
	} else {
		//兼容模式
		log.Println("兼容模式")
		f.Call([]reflect.Value{reflect.ValueOf(request), reflect.ValueOf(response)})
	}

	response.OriginResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if crossDomain != "" {
		response.OriginResponseWriter.Header().Set("Access-Control-Allow-Origin", crossDomain)
	}

}
