package users

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"server-monitoring/utils/rest_error"
	"strconv"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Id             int64              `bson:"id"`
	FirstName      string             `bson:"first_name"`
	UserName       string             `bson:"username"`
	Password       string             `bson:"password"`
	Email          string             `bson:"email"`
	LastName       string             `bson:"last_name"`
	RealName       string             `bson:"real_name"`
	ProfileImageId string             `bson:"profile_image_id"`
	Status         int8               `bson:"status"`
	Gender         int8               `bson:"gender"`
	IsTeacher      int8               `bson:"is_teacher"`
	IsSuperAdmin   int8               `bson:"is_super_admin"`
	PhoneNumber    string             `bson:"phone_number"`
	VerifyCode     string             `bson:"verify_code"`
	CreatedAt      string             `bson:"created_at"`
	UpdatedAt      string             `bson:"updated_at"`
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
