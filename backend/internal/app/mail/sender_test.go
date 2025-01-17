package mail

import (
	"strings"
	"testing"
)

func TestSendToSupportWithFile(t *testing.T) {
	sender := NewSender("sabdpp17@gmail.com", "secret")
	fileBASE64 := strings.Repeat("0", 1024*1024)
	err := sender.SendToSupportWithFile("test", "test body", "sabdpp17@gmail.com", "lol.png", &fileBASE64)
	if err == nil {
		t.Error("cannot handle big file", err)
	}
}
