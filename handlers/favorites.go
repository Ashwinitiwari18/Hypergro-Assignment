package handlers

import (
	"context"
	"net/http"
	"time"

	"property-listing/db"
	"property-listing/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecommendRequest struct {
	ToUserEmail string `json:"toUserEmail" binding:"required,email"`
	Message     string `json:"message"`
}

func AddFavorite(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))
	propertyID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	// Check if property exists
	var property models.Property
	err = db.Properties.FindOne(context.Background(), bson.M{"_id": propertyID}).Decode(&property)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Check if already favorited
	var existingFavorite models.Favorite
	err = db.Favorites.FindOne(context.Background(), bson.M{
		"userId":     userID,
		"propertyId": propertyID,
	}).Decode(&existingFavorite)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Property already in favorites"})
		return
	}

	favorite := models.Favorite{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		PropertyID: propertyID,
		CreatedAt:  time.Now(),
	}

	_, err = db.Favorites.InsertOne(context.Background(), favorite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding to favorites"})
		return
	}

	c.JSON(http.StatusCreated, favorite)
}

func RemoveFavorite(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))
	propertyID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	result, err := db.Favorites.DeleteOne(context.Background(), bson.M{
		"userId":     userID,
		"propertyId": propertyID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing from favorites"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

func GetFavorites(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))

	// Get all favorites for user
	cursor, err := db.Favorites.Find(context.Background(), bson.M{"userId": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching favorites"})
		return
	}
	defer cursor.Close(context.Background())

	var favorites []models.Favorite
	if err = cursor.All(context.Background(), &favorites); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding favorites"})
		return
	}

	// Get property details for each favorite
	var properties []models.Property
	for _, fav := range favorites {
		var property models.Property
		err := db.Properties.FindOne(context.Background(), bson.M{"_id": fav.PropertyID}).Decode(&property)
		if err == nil {
			properties = append(properties, property)
		}
	}

	c.JSON(http.StatusOK, properties)
}

func RecommendProperty(c *gin.Context) {
	fromUserID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))
	propertyID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	var req RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find recipient user
	var toUser models.User
	err = db.Users.FindOne(context.Background(), bson.M{"email": req.ToUserEmail}).Decode(&toUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient user not found"})
		return
	}

	// Check if property exists
	var property models.Property
	err = db.Properties.FindOne(context.Background(), bson.M{"_id": propertyID}).Decode(&property)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	recommendation := models.Recommendation{
		ID:         primitive.NewObjectID(),
		FromUserID: fromUserID,
		ToUserID:   toUser.ID,
		PropertyID: propertyID,
		Message:    req.Message,
		CreatedAt:  time.Now(),
		IsRead:     false,
	}

	_, err = db.Recommendations.InsertOne(context.Background(), recommendation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating recommendation"})
		return
	}

	c.JSON(http.StatusCreated, recommendation)
}

func GetRecommendations(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))

	// Get all recommendations for user
	cursor, err := db.Recommendations.Find(context.Background(), bson.M{"toUserId": userID},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommendations"})
		return
	}
	defer cursor.Close(context.Background())

	var recommendations []models.Recommendation
	if err = cursor.All(context.Background(), &recommendations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding recommendations"})
		return
	}

	// Get property and sender details for each recommendation
	type RecommendationWithDetails struct {
		models.Recommendation
		Property models.Property `json:"property"`
		Sender   models.User     `json:"sender"`
	}

	var detailedRecommendations []RecommendationWithDetails
	for _, rec := range recommendations {
		var property models.Property
		var sender models.User

		err := db.Properties.FindOne(context.Background(), bson.M{"_id": rec.PropertyID}).Decode(&property)
		if err != nil {
			continue
		}

		err = db.Users.FindOne(context.Background(), bson.M{"_id": rec.FromUserID}).Decode(&sender)
		if err != nil {
			continue
		}

		detailedRecommendations = append(detailedRecommendations, RecommendationWithDetails{
			Recommendation: rec,
			Property:       property,
			Sender:         sender,
		})
	}

	c.JSON(http.StatusOK, detailedRecommendations)
}

func MarkRecommendationAsRead(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))
	recommendationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recommendation ID"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"isRead": true,
		},
	}

	_, err = db.Recommendations.UpdateOne(context.Background(), bson.M{
		"_id":      recommendationID,
		"toUserId": userID,
	}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marking recommendation as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recommendation marked as read"})
}
