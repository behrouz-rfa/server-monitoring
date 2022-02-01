package settings

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"server-monitoring/shared/database"
)

func (l *Setting) Update() error {
	coll := database.Mongo.Database("monitoring").Collection("settings")
	setting := bson.D{
		{"site_name", l.SiteName},
		{"language_name", l.LanguageName},
		{"language_id", l.LanguageId},
		{"meta", l.Meta},
		{"keyword", l.Keyword},
		{"email", l.Email},
		{"password", l.Password},
		{"tel", l.Tel},
		{"phone", l.Phone},
		{"username", l.Username},
		{"interface", l.Interface},
		{"filter", l.Filter},
		{"status", l.Status},
		{"message", l.Message},
	}
	err := coll.FindOneAndReplace(context.Background(), bson.D{{}}, setting).Decode(&l)
	if err != nil {
		return err
	}
	return nil
}

func (l *Setting) FindFist() error {
	coll := database.Mongo.Database("monitoring").Collection("settings")
	if err := coll.FindOne(context.Background(), bson.M{}).Decode(&l); err != nil {
		return err
	}

	return nil
}

func (l *Setting) FindByUserName() error {
	coll := database.Mongo.Database("monitoring").Collection("settings")
	if err := coll.FindOne(context.Background(), bson.M{"username": l.Username}).Decode(&l); err != nil {
		return err
	}

	return nil
}
