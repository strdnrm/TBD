package model_test

import (
	"final_project/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Valdiate(t *testing.T) {
	// u := model.TestUser(t)
	// assert.NoError(t, u.Validate())
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = "adadsa.gmail.com"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "012"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
