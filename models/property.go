package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title         string             `json:"title" bson:"title"`
	Type          string             `json:"type" bson:"type"`
	Price         float64            `json:"price" bson:"price"`
	State         string             `json:"state" bson:"state"`
	City          string             `json:"city" bson:"city"`
	AreaSqFt      float64            `json:"areaSqFt" bson:"areaSqFt"`
	Bedrooms      int                `json:"bedrooms" bson:"bedrooms"`
	Bathrooms     int                `json:"bathrooms" bson:"bathrooms"`
	Amenities     []string           `json:"amenities" bson:"amenities"`
	Furnished     string             `json:"furnished" bson:"furnished"`
	AvailableFrom string             `json:"availableFrom" bson:"availableFrom"`
	ListedBy      string             `json:"listedBy" bson:"listedBy"`
	Tags          []string           `json:"tags" bson:"tags"`
	ColorTheme    string             `json:"colorTheme" bson:"colorTheme"`
	Rating        float64            `json:"rating" bson:"rating"`
	IsVerified    bool               `json:"isVerified" bson:"isVerified"`
	ListingType   string             `json:"listingType" bson:"listingType"`
	CreatedBy     primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}
