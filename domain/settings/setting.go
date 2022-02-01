package settings

import "go.mongodb.org/mongo-driver/bson/primitive"

type Setting struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SiteName     string             `bson:"site_name"`
	LanguageName string             `bson:"language_name"`
	IsSuperAdmin int                `bson:"is_super_admin"`
	LanguageId   int                `bson:"language_id"`
	Meta         string             `bson:"meta"`
	Keyword      string             `bson:"keyword"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	Tel          string             `bson:"tel"`
	Phone        string             `bson:"phone"`
	Username     string             `bson:"username"`
	Interface    string             `bson:"interface"`
	Filter       string             `bson:"filter"`
	Status       int8               `bson:"status"`
	Message      string             `bson:"message"`
}
