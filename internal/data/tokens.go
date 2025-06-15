package data

import (
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/validator"
)

// TODO: make this into Enum
const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plain  string    `json:"token"`
	Hash   []byte    `json:"-"`
	UserID int64     `json:"-"`
	Expiry time.Time `json:"expiry"`
	Scope  string    `json:"-"`
}

func ValidatePlaintextToken(v *validator.Validator, token string) {
	v.Check(token != "", "token", "token must be provided")
	v.Check(len(token) < 32, "token", "must not be less larger than 32 bytes long")
}
