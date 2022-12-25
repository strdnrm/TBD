package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Login:    "user",
		Email:    "user@example.com",
		Password: "Password",
		Name:     "user",
		Surname:  "usersur",
		Is_admin: false,
	}
}
