package model

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                uuid.UUID `json:"-" db:"id" validate:"uuid"`
	Login             string    `json:"login" db:"login" validate:"required"`
	Email             string    `json:"email" db:"email" validate:"required,email"`
	Password          string    `json:"password,omitempty" db:"-" validate:"required,gte=6,lte=40"`
	EncryptedPassword string    `json:"-" db:"password"`
	Name              string    `json:"name" db:"name" validate:"required"`
	Surname           string    `json:"surname" db:"surname" validate:"required"`
	Is_admin          bool      `json:"-" db:"is_admin" `
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
		fmt.Println(enc)
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

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
