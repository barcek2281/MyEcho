package model

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Date     sql.NullTime
}

func (u *Admin) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (a *Admin) BeforeCreate() error {
	if len(a.Password) > 0 {
		enc, err := Encrypt(a.Password)
		if err != nil {
			return err
		}
		a.Password = enc
	}
	return nil
}
