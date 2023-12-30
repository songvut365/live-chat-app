package model

import "time"

type Message struct {
	Id        string
	Sender    string
	Content   string
	Timestamp time.Time
}
