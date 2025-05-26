package utils

import (
	"context"
	"log"

	"property-listing/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateCreatedByForAllProperties sets createdBy for all properties to the given ObjectID
func UpdateCreatedByForAllProperties() error {
	ctx := context.Background()
	newID, err := primitive.ObjectIDFromHex("6834cabbc0944cc269a10ec6")
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"createdBy": newID}}
	_, err = db.Properties.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		log.Printf("Error updating createdBy for all properties: %v", err)
		return err
	}
	log.Println("Successfully updated createdBy for all properties.")
	return nil
}
