package helper

import "github.com/satori/go.uuid"

func CreateUUID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}
