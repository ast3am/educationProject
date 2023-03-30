package models

type UserModel struct {
	ID      string       `json:"id" bson:"id"`
	Name    string       `json:"name" bson:"name"`
	Age     string       `json:"age" bson:"age"`
	Friends []*UserModel `json:"friends" bson:"-"`
}
