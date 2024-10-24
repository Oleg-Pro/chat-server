package model

// ChatInfo info about chat
type ChatInfo struct {
	Users string
}

// Chat entity
type Chat struct {
	ID   int64
	Info ChatInfo
}
