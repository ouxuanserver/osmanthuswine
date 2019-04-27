package helper

import "testing"

func TestCreateUUID(t *testing.T) {
	s := CreateUUID()
	println(s)
	if len(s) != 36 {
		t.Error("UUID长度错误")
	}
}
