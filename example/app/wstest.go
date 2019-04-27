package app

import (
	"github.com/wailovet/osmanthuswine/src/core"
	"gopkg.in/olahol/melody.v1"
)

type Wstest struct {
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
