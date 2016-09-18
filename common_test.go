package appcomm

import (
	"testing"
)

func TestMd5(t *testing.T) {
	s := Md5("aaaaa")
	t.Logf("%s\n", s)
}
func TestDateString(t *testing.T) {
	s := GetDateString()
	t.Logf("%s", s)
}
