package models

import "time"

type Chat struct {
	Id         int
	Name       string
	LastActive time.Time
}
