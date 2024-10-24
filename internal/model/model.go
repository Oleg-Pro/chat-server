package model

import (
	"database/sql"
)

// ChatInfo info about chat
type ChatInfo struct {
	Users string
}

// Chat entity
type Chat struct {
	ID   int64
	Info ChatInfo
}

type MessageInfo struct {
	From      string
	Text      string
	Timestamp sql.NullTime
}
