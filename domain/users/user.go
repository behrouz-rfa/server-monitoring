package users

import (
	"errors"
	"math/rand"
	"server-monitoring/utils/rest_error"
	"strconv"
)

type User struct {
	Id             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	UserName       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	LastName       string `json:"last_name"`
	RealName       string `json:"real_name"`
	ProfileImageId string `json:"profile_image_id"`
	Status         int8   `json:"status"`
	Gender         int8   `json:"gender"`
	IsTeacher      int8   `json:"is_teacher"`
	IsSuperAdmin   int8   `json:"is_super_admin"`
	PhoneNumber    string `json:"phone_number"`
	VerifyCode     string `json:"verify_code"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func (u *User) Validate() rest_error.RestErr {
	if len(u.PhoneNumber) < 10 {
		return rest_error.NewInternalServerError("phone number invaide", errors.New("invalide"))
	}
	u.VerifyCode = strconv.Itoa(rand.Intn(9999-1000) + 1000)
	return nil
}

func (u *User) ValidatePhone() rest_error.RestErr {
	if len(u.PhoneNumber) < 10 {
		return rest_error.NewInternalServerError("phone number is invalid", errors.New("invalide"))
	}
	return nil
}

