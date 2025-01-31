package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tc := []struct {
		name    string
		user    User
		isValid bool
	}{
		{
			name:    "Validate OK",
			user:    createuser(1, "sabdpp12@gma.ca", "barcek2281", "1234567"),
			isValid: true,
		},
		{
			name:    "Incorrect Password",
			user:    createuser(1, "sabdpp12@gma.ca", "barcek2281", "12"),
			isValid: false,
		},
		{
			name:    "Incorrect Password",
			user:    createuser(1, "sabdpp12@gma.ca", "barcek2281", ""),
			isValid: false,
		},
	}

	for _, testik := range tc {
		t.Run(testik.name, func(t *testing.T) {
			if testik.isValid {
				assert.NoError(t, testik.user.Validate())
			} else {
				assert.Error(t, testik.user.Validate())
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	u := createuser(1, "", "", "123454")
	err := u.BeforeCreate()
	assert.Nil(t, err)

	if !u.ComparePassword("123454") {
		t.Error("password should hash and compair password", u.ComparePassword("123454"))
	}

	if u.ComparePassword("12") {
		t.Error("Didnt mantch with correct password")
	}
}

func TestEncryptPassword(t *testing.T) {
	password := strings.Repeat("0", 10)
	_, err := Encrypt(password)
	assert.Nil(t, err)

	password = strings.Repeat("0", 100)
	_, err = Encrypt(password)
	assert.NotNil(t, err)
}

func createuser(id int, email, login, password string) User {
	return User{
		ID:       id,
		Email:    email,
		Login:    login,
		Password: password,
	}
}
