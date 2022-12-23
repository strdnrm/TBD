package model

import (
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int    `json:"-" db:"id"` //validate:"required, numeric"
	Login             string `json:"login" db:"login" validate:"required"`
	Email             string `json:"email" db:"email" validate:"required,email"`
	Password          string `json:"password,omitempty" db:"password" validate:"required,gte=6,lte=40"`
	EncryptedPassword string `json:"-"` //`db:"id" validate:"required"`
	Name              string `json:"name" db:"name" validate:"required"`
	Surname           string `json:"surname" db:"surname" validate:"required"`
	Is_admin          bool   `json:"-" db:"is_admin" ` //validate:"required"
}

func (u *User) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
