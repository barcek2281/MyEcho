package model

import (
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

func createuser(id int, email, login, password string) User {
	return User{
		ID:       id,
		Email:    email,
		Login:    login,
		Password: password,
	}
}
