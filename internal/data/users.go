package data

import (
	"errors"
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"-"`
}

func (u *User) IsAnonymousUser() bool {
	return u == AnonymousUser
}

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlainText(v *validator.Validator, platintextPassword string) {
	v.Check(platintextPassword != "", "password", "must be provided")
	v.Check(len(platintextPassword) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(platintextPassword) <= 72, "password", "must not be 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlainText(v, *user.Password.Plaintext)
	}
	// Should panic here because if this happens, we forget to set the hash password
	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
