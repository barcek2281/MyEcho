package model

import "time"

type Messages struct {
	Id       int
	Sender   string
	Receiver string
	Message  string
	Date     time.Time
}
