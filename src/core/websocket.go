package core

import (
	"github.com/wailovet/osmanthuswine/src/interfaces"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"sync"
)

var instanceMelody sync.Map

type WebSocket struct {
	ws      *melody.Melody
	FunName string
}

func (that *WebSocket) SetFunName(funName string) {
	that.FunName = funName
}

func (that *WebSocket) WebSocketInit(wsx *melody.Melody) {
	that.ws = wsx
}
func (that *WebSocket) GetMelody() *melody.Melody {
	return that.ws
}

func (that *WebSocket) GetWebSocket() *melody.Melody {
	return that.ws
}

func GetWebSocket(group string, hand interfaces.WebSocketInterface) *melody.Melody {
	var m *melody.Melody
	tmp, ok := instanceMelody.Load(group)
	if !ok {
		m = melody.New()
		m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		m.HandleConnect(hand.HandleConnect)
		m.HandleDisconnect(hand.HandleDisconnect)
		m.HandlePong(hand.HandlePong)
		m.HandleError(hand.HandleError)
		m.HandleMessage(hand.HandleMessage)
		m.HandleMessageBinary(hand.HandleMessageBinary)
		m.HandleSentMessage(hand.HandleSentMessage)
		m.HandleSentMessageBinary(hand.HandleSentMessageBinary)
		instanceMelody.Store(group, m)
	} else {
		m = tmp.(*melody.Melody)
	}
	return m
}
