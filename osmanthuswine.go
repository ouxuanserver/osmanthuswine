package osmanthuswine

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ouxuanserver/osmanthuswine/src/core"
	"github.com/ouxuanserver/osmanthuswine/src/helper"
	"github.com/ouxuanserver/osmanthuswine/src/session"
	"github.com/wailovet/overseer"
	"github.com/wailovet/overseer/fetcher"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var chiRouter *chi.Mux

func init() {
	//获取相对于执行文件的工作目录的绝对路径，并且把路径设置为工作目录
	if err := os.Chdir(filepath.Dir(os.Args[0])); err != nil {
		log.Fatal("设置工作目录失败：", err)
	}
}
func GetChiRouter() *chi.Mux {
	if chiRouter == nil {

		chiRouter = chi.NewRouter()
		chiRouter.Use(middleware.RequestID)
		chiRouter.Use(middleware.RealIP)
		chiRouter.Use(middleware.Logger)
		chiRouter.Use(middleware.Recoverer)
		chiRouter.Use(middleware.Timeout(60 * time.Second))
	}
	return chiRouter
}
func Run() {
	path, _ := os.Getwd()

	log.Println("工作目录:", path)
	cc := core.GetInstanceConfig()

	if runtime.GOOS == "windows" || cc.UpdatePath == "" {
		listener, err := net.Listen("tcp", cc.Host+":"+cc.Port)
		if err != nil {
			log.Fatal(err.Error())
		}
		RunProg(overseer.State{
			Listener: listener,
		})
	} else {
		overseer.Run(overseer.Config{
			Program: RunProg,
			Address: cc.Host + ":" + cc.Port,
			Fetcher: &fetcher.File{
				Path:     cc.UpdateDir + cc.UpdatePath,
				Interval: time.Second * 10,
			},
		})
	}

}
func RunProg(state overseer.State) {

	cc := core.GetInstanceConfig()

	helper.GetInstanceLog().Out("开始监听:", cc.Host+":"+cc.Port)

	r := GetChiRouter()

	apiRouter := cc.ApiRouter

	r.HandleFunc(apiRouter, func(writer http.ResponseWriter, request *http.Request) {

		requestData := core.Request{}

		sessionMan := session.New(request, writer)

		requestData.REQUEST = make(map[string]string)
		//GET
		requestData.SyncGetData(request)
		//POST
		requestData.SyncPostData(request, cc.PostMaxMemory)
		//HEADER
		requestData.SyncHeaderData(request)
		//COOKIE
		requestData.SyncCookieData(request)
		//SESSION
		requestData.SyncSessionData(sessionMan)

		responseHandle := core.Response{OriginResponseWriter: writer, Session: sessionMan}

		defer func() {
			errs := recover()
			if errs == nil {
				return
			}
			errtxt := fmt.Sprintf("%v", errs)
			if errtxt != "" {
				responseHandle.DisplayByError(errtxt, 500, strings.Split(string(debug.Stack()), "\n\t")...)
			}
		}()

		core.GetInstanceRouterManage().RouterSend(request.URL.Path, requestData, responseHandle, cc.CrossDomain)

	})

	r.Handle("/*", http.FileServer(http.Dir("html")))
	//r.HandleFunc("/html/*", func(writer http.ResponseWriter, request *http.Request) {
	//	path := request.URL.Path
	//	if path == "/html/" {
	//		path = "/index.html"
	//	}
	//
	//	path=strings.TrimLeft(path,"/")
	//	helper.GetInstanceLog().Out("静态文件:", path)
	//
	//	f, err := os.Stat(path)
	//	if err == nil {
	//		if f.IsDir() {
	//			path += "/index.html"
	//		}
	//		data, err := ioutil.ReadFile(path)
	//		if err == nil {
	//			writer.Write(data)
	//			return
	//		}
	//	}
	//
	//	writer.WriteHeader(404)
	//	writer.Write([]byte(err.Error()))
	//
	//})

	if err := http.Serve(state.Listener, r); err != nil {
		log.Fatal(err)
	}
	//http.ListenAndServe(cc.Host+":"+cc.Port, r)

}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}
