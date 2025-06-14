package data

import (
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/validator"
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

func ValidatePlaintextToken(v *validator.Validator, token string) {
	v.Check(token != "", "token", "token must be provided")
	v.Check(len(token) < 32, "token", "must not be less larger than 32 bytes long")
}
