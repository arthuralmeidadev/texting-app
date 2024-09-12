package models

import "time"

type FriendRequest struct {
    SendUsrn string
    RecUsrn string
    Status string
    LastMod time.Time
}
