package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Login, validation.Required, validation.NilOrNotEmpty),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 71)),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := Encrypt(u.Password)
		if err != nil {
			return err
		}
		u.Password = enc
	}
	return nil
}

func Encrypt(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
