## 内部使用基于go-chi的web框架

# 框架引入
> go get -u github.com/wailovet/osmanthuswine

### 目录结构
```
app
 |--index
      |--index.go
html
 |--静态文件....
main.go
config.json
```

# 开始
#### 创建以上目录结构


+ /config.json 配置文件

```json
{
  "port": "8808",
  "host": "0.0.0.0",
  "cross_domain": "*",
  "post_max_memory": 1024000,
  "update_path": "new_exe",
  "db": {
    "host": "",
    "port": "",
    "user": "",
    "password": "",
    "name": "",
    "max_open_conn": 500
  }
}
```

+ /main.go文件

```
package main

import (
	"./app/index"
	"github.com/wailovet/osmanthuswine"
	"github.com/wailovet/osmanthuswine/src/core"
)

func main() {
	//注册index控制器
	core.GetInstanceRouterManage().Registered(&index.Index{})
	//主程序执行
	osmanthuswine.Run()
}
```


+ /app/index/index.go文件

```
package index

import (
    "github.com/wailovet/osmanthuswine/src/core"
)

type Index struct {
    core.Controller
}

func (that *Index) Index() {
    that.DisplayByData(that.Request.REQUEST)
}

```
## core.Controller.Request 使用说明
#### 输入参数 
+ that.Request.GET["参数"] //类型:map[string]string
+ that.Request.POST["参数"] //类型:map[string]string
+ that.Request.REQUEST["参数"] //类型:map[string]string , 为GET以及POST的合并值,当出现值冲突时GET参数会被覆盖
#### session与cookie获取
+ that.Request.SESSION["参数"] //类型:map[string]string
+ that.Request.COOKIE["参数"] //类型:map[string]string
#### header信息获取
+ that.Request.HEADER["参数"] //类型:map[string]string
#### 该取值一般以POST-RAW形式传入原始数据,有可能
+ that.Request.BODY //类型:string
#### 快速获取上传的文件
+ that.Request.FILE //类型:*multipart.FileHeader
#### 获取上传的所有文件
+ that.Request.FILES //类型:map[string][]*multipart.FileHeader



## core.Controller 使用说明
#### 输出显示
> 即使输出时不处于函数结尾,也无需return
+ that.DisplayByData(data interface{})
```
{
    "code":0,
    "data":data,
    "msg":""
}
```
- - -
+ that.DisplayBySuccess(msg string)
```
{
    "code":0,
    "data":null,
    "msg":msg
}
```
- - -
+ that.DisplayByError(msg string, code int)
```
{
    "code":code,
    "data":null,
    "msg":msg
}
```
- - -
+ that.Display(data interface{}, msg string, code int)
```
{
    "code":code,
    "data":data,
    "msg":msg
}
```
- - -
+ that.DisplayByString(data string)
```
data //直接输出data以string形式
```
- - -
+ that.DisplayByRaw(data []byte)
```
data //直接输出data以[]byte形式,可用于直接输出二进制文件
```
- - -
##### 20190416新增
+ that.CheckErrDisplayByError(err error,msg...)
```
err //错误信息,自动判断是否等于nil,如果等于nil该语句会被忽略
msg //错误文案提示,不填直接输出err.Error()
```

#### session操作
> 目前session实现基于securecookie,以加密形式储存在cookie中,注意不要存放大量数据,以免超过cookie的最大储存值
+ that.SetSession(name string, value string) //设置session
+ that.DeleteSession(name string) //删除session
+ that.ClearSession() //清空session

#### cookie操作
> 尽量以session的形式操作
+ that.SetCookie(name string, value string)


## 数据库操作
> 目前框架中集成gorm与xorm框架
+ core.GetXormAuto() //获取xorm实例
+ core.GetGormAuto() //获取gorm实例
#### 数据库配置
```
实例的数据库配置来自于相同目录下的config.json或者private.json文件
{
  ...其他配置
  "db": {
    "host": "",
    "port": "",
    "user": "",
    "password": "",
    "name": "",
    "prefix": "",
    "max_open_conn": 500
  }
} 
prefix为表前缀
max_open_conn为可支持最大连接数(未测试是否可用
```

## 支持WebSocket
> 当传入core.GetInstanceRouterManage().Registered的对象继承自core.WebSocket时,协议升级为websocket,路由地址忽略最后方法名


#### 集成melody库,使用详情https://github.com/olahol/melody
```
package index

import (
	"github.com/wailovet/osmanthuswine/src/core"
	"gopkg.in/olahol/melody.v1"
)

type Index struct {
	core.WebSocket
}

func (that *Wstest) HandleConnect(session *melody.Session) {
	//implement
}

func (that *Wstest) HandlePong(session *melody.Session) {
	//implement
}

func (that *Wstest) HandleMessage(session *melody.Session, data []byte) {
	that.GetMelody().Broadcast(data)
	//implement
}

func (that *Wstest) HandleMessageBinary(session *melody.Session, data []byte) {
	//implement
}

func (that *Wstest) HandleSentMessage(session *melody.Session, data []byte) {
	//implement
}

func (that *Wstest) HandleSentMessageBinary(session *melody.Session, data []byte) {
	//implement
}

func (that *Wstest) HandleDisconnect(session *melody.Session) {
	//implement
}

func (that *Wstest) HandleError(session *melody.Session, err error) {
	//implement
}

```

```
//javascript
var ws = new WebSocket("ws://127.0.0.1/Api/Index/Index")
```
> PS:不同url对应不同的melody实例


## 杂项
> 热更新,仅支持linux
```
默认情况下,检测同路径下的<文件名_update>,如果该文件与当前文件不一致,则进行热更,已连接的连接无需断连
可在config.json中配置检测的文件名
{
    ...其他配置
    "update_path": "需要检测的文件路径"
}
备注:需要检测的文件路径最好不要与当前运行的文件路径相同
```