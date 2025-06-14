package data

import (
	"time"
)

// TODO: make this into Enum
const (
	ScopeActivation = "activation"
)

type Token struct {
	Plain  string
	Hash   []byte
	UserID int64
	Expiry time.Time
	Scope  string
}
