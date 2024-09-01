package models

import "time"

type Message struct {
	Id        int
	Sender    string
	Chat      int
	Content   string
	SentAt    time.Time
	RepliesTo int
}
