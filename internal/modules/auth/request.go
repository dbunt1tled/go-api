package auth

import (
	"time"

	"github.com/dbunt1tled/go-api/internal/modules/user"
)

type Login struct {
	Email    string `json:"email"    validate:"required,email,max=70"        example:"fake@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20,passwd" example:"pas$word1A"`
}

type Register struct {
	FirstName       string `json:"firstName"       validate:"required,min=2"                                       example:"John"`
	SecondName      string `json:"secondName"      validate:"required,min=2"                                       example:"Dou"`
	Email           string `json:"email"           validate:"required,email,unique_db=users.email"                 example:"example@example.com"`
	PhoneNumber     string `json:"phoneNumber"     validate:"required,unique_db=users.phone_number"                example:"+1234567890"`
	Password        string `json:"password"        validate:"required,min=8,max=20,passwd,eqfield=PasswordConfirm" example:"pas$word1A"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"                                             example:"pas$word1A"`
}

func (r Register) ToUser() *user.User {
	return (&user.User{
		FirstName:   r.FirstName,
		SecondName:  r.SecondName,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
		Status:      user.Pending,
		Roles:       []user.Role{user.Client},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}).NextID()
}

type Confirm struct {
	Token string `params:"token" json:"token" validate:"required,min=8" example:"random string"`
}
