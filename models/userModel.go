package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID				primitive.ObjectID		`bson:"_id"`
	First_Name		*string					`json:"first_name" validate:"required,min=2,max=100"`
	Last_Name		*string					`json:"last_name" validate:"required,min=2,max=100"`
	Password		*string					`json:"Password" validate:"required,min=6"`
	Email			*string					`json:"email" validate:"email,required"`
	User_Type		*string					`json:"user_type" validate:"required,eq=USER|eq=ADMIN"`
	Avatar			*string					`json:"avatar"`
	Phone			*string					`json:"phone" validate:"required"`
	Token			*string					`json:"token"`
	Refresh_Token	*string					`json:"refresh_token"`
	Created_At		time.Time				`json:"created_at"`
	Updated_At		time.Time				`json:"updated_at"`
	User_id			string					`json:"user_id"`
}