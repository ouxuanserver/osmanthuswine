package interfaces

import (
	"gopkg.in/olahol/melody.v1"
)

//WebSocketInterface interface
type WebSocketInterface interface {
	SetFunName(funName string)
	WebSocketInit(ws *melody.Melody)
	GetWebSocket() *melody.Melody
	GetMelody() *melody.Melody
	HandleConnect(*melody.Session)
	HandlePong(*melody.Session)
	HandleMessage(*melody.Session, []byte)
	HandleMessageBinary(*melody.Session, []byte)
	HandleSentMessage(*melody.Session, []byte)
	HandleSentMessageBinary(*melody.Session, []byte)
	HandleDisconnect(*melody.Session)
	HandleError(*melody.Session, error)
}
