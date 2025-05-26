package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string            `json:"email" bson:"email"`
	Password  string            `json:"-" bson:"password"`
	Name      string            `json:"name" bson:"name"`
	CreatedAt time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type Favorite struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	PropertyID primitive.ObjectID `json:"propertyId" bson:"propertyId"`
	CreatedAt  time.Time         `json:"createdAt" bson:"createdAt"`
}

type Recommendation struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FromUserID   primitive.ObjectID `json:"fromUserId" bson:"fromUserId"`
	ToUserID     primitive.ObjectID `json:"toUserId" bson:"toUserId"`
	PropertyID   primitive.ObjectID `json:"propertyId" bson:"propertyId"`
	Message      string            `json:"message" bson:"message"`
	CreatedAt    time.Time         `json:"createdAt" bson:"createdAt"`
	IsRead       bool              `json:"isRead" bson:"isRead"`
} 