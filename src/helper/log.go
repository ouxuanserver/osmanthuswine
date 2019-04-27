package helper

import "log"

type Log struct {
}

var instanceLog *Log

func GetInstanceLog() *Log {
	if instanceLog == nil {
		instanceLog = &Log{} // not thread safe
	}
	return instanceLog
}

func (l *Log) Out(args ...interface{}) {
	log.Println(args)
}
