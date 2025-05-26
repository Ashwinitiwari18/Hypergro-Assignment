package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"property-listing/cache"
	"property-listing/db"
	"property-listing/models"

	"crypto/sha1"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropertyRequest struct {
	Title         string   `json:"title" binding:"required"`
	Type          string   `json:"type" binding:"required"`
	Price         float64  `json:"price" binding:"required"`
	State         string   `json:"state" binding:"required"`
	City          string   `json:"city" binding:"required"`
	AreaSqFt      float64  `json:"areaSqFt" binding:"required"`
	Bedrooms      int      `json:"bedrooms" binding:"required"`
	Bathrooms     int      `json:"bathrooms" binding:"required"`
	Amenities     []string `json:"amenities"`
	Furnished     string   `json:"furnished"`
	AvailableFrom string   `json:"availableFrom"`
	ListedBy      string   `json:"listedBy"`
	Tags          []string `json:"tags"`
	ColorTheme    string   `json:"colorTheme"`
	Rating        float64  `json:"rating"`
	IsVerified    bool     `json:"isVerified"`
	ListingType   string   `json:"listingType"`
	Location      string   `json:"location"`
	Area          float64  `json:"area"`
	Features      []string `json:"features"`
	Status        string   `json:"status"`
	Description   string   `json:"description"`
}

func CreateProperty(c *gin.Context) {
	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))
	property := models.Property{
		ID:            primitive.NewObjectID(),
		Title:         req.Title,
		Type:          req.Type,
		Price:         req.Price,
		State:         req.State,
		City:          req.City,
		AreaSqFt:      req.AreaSqFt,
		Bedrooms:      req.Bedrooms,
		Bathrooms:     req.Bathrooms,
		Amenities:     req.Amenities,
		Furnished:     req.Furnished,
		AvailableFrom: req.AvailableFrom,
		ListedBy:      req.ListedBy,
		Tags:          req.Tags,
		ColorTheme:    req.ColorTheme,
		Rating:        req.Rating,
		IsVerified:    req.IsVerified,
		ListingType:   req.ListingType,
		Location:      req.Location,
		Area:          req.Area,
		Features:      req.Features,
		Status:        req.Status,
		Description:   req.Description,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := db.Properties.InsertOne(context.Background(), property)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating property"})
		return
	}

	deleteAllPropertyListCaches()
	c.JSON(http.StatusCreated, property)
}

func GetProperties(c *gin.Context) {
	// Build a cache key based on the query string
	queryString := c.Request.URL.RawQuery
	hash := sha1.Sum([]byte(queryString))
	cacheKey := "properties:list:" + hex.EncodeToString(hash[:])

	var properties []models.Property
	err := cache.Get(cacheKey, &properties)
	if err == nil {
		c.JSON(http.StatusOK, properties)
		return
	}

	// Build filter from query parameters
	filter := bson.M{}
	if priceMin := c.Query("priceMin"); priceMin != "" {
		price, _ := strconv.ParseFloat(priceMin, 64)
		filter["price"] = bson.M{"$gte": price}
	}
	if priceMax := c.Query("priceMax"); priceMax != "" {
		price, _ := strconv.ParseFloat(priceMax, 64)
		if filter["price"] != nil {
			filter["price"].(bson.M)["$lte"] = price
		} else {
			filter["price"] = bson.M{"$lte": price}
		}
	}
	if location := c.Query("location"); location != "" {
		filter["location"] = bson.M{"$regex": location, "$options": "i"}
	}
	if bedrooms := c.Query("bedrooms"); bedrooms != "" {
		bed, _ := strconv.Atoi(bedrooms)
		filter["bedrooms"] = bed
	}
	if propertyType := c.Query("type"); propertyType != "" {
		filter["type"] = propertyType
	}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}
	if state := c.Query("state"); state != "" {
		filter["state"] = state
	}
	if city := c.Query("city"); city != "" {
		filter["city"] = city
	}
	if areaSqFtMin := c.Query("areaSqFtMin"); areaSqFtMin != "" {
		area, _ := strconv.ParseFloat(areaSqFtMin, 64)
		filter["areaSqFt"] = bson.M{"$gte": area}
	}
	if areaSqFtMax := c.Query("areaSqFtMax"); areaSqFtMax != "" {
		area, _ := strconv.ParseFloat(areaSqFtMax, 64)
		if filter["areaSqFt"] != nil {
			filter["areaSqFt"].(bson.M)["$lte"] = area
		} else {
			filter["areaSqFt"] = bson.M{"$lte": area}
		}
	}
	if furnished := c.Query("furnished"); furnished != "" {
		filter["furnished"] = furnished
	}
	if availableFrom := c.Query("availableFrom"); availableFrom != "" {
		filter["availableFrom"] = availableFrom
	}
	if listedBy := c.Query("listedBy"); listedBy != "" {
		filter["listedBy"] = listedBy
	}
	if colorTheme := c.Query("colorTheme"); colorTheme != "" {
		filter["colorTheme"] = colorTheme
	}
	if ratingMin := c.Query("ratingMin"); ratingMin != "" {
		rating, _ := strconv.ParseFloat(ratingMin, 64)
		filter["rating"] = bson.M{"$gte": rating}
	}
	if ratingMax := c.Query("ratingMax"); ratingMax != "" {
		rating, _ := strconv.ParseFloat(ratingMax, 64)
		if filter["rating"] != nil {
			filter["rating"].(bson.M)["$lte"] = rating
		} else {
			filter["rating"] = bson.M{"$lte": rating}
		}
	}
	if isVerified := c.Query("isVerified"); isVerified != "" {
		if isVerified == "true" || isVerified == "1" {
			filter["isVerified"] = true
		} else if isVerified == "false" || isVerified == "0" {
			filter["isVerified"] = false
		}
	}
	if listingType := c.Query("listingType"); listingType != "" {
		filter["listingType"] = listingType
	}
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter["tags"] = bson.M{"$in": tags}
	}
	if amenities := c.QueryArray("amenities"); len(amenities) > 0 {
		filter["amenities"] = bson.M{"$all": amenities}
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := db.Properties.Find(context.Background(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching properties"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &properties); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding properties"})
		return
	}

	// Cache the results for this filter
	cache.Set(cacheKey, properties, time.Hour)

	c.JSON(http.StatusOK, properties)
}

func GetProperty(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	// Try to get from cache first
	var property models.Property
	err = cache.Get("property:"+id, &property)
	if err == nil {
		c.JSON(http.StatusOK, property)
		return
	}

	err = db.Properties.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&property)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	// Cache the result
	cache.Set("property:"+id, property, time.Hour)

	c.JSON(http.StatusOK, property)
}

func UpdateProperty(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))

	// Check if property exists and belongs to user
	var property models.Property
	err = db.Properties.FindOne(context.Background(), bson.M{
		"_id":       objectID,
		"createdBy": userID,
	}).Decode(&property)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found or unauthorized"})
		return
	}

	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":         req.Title,
			"type":          req.Type,
			"price":         req.Price,
			"state":         req.State,
			"city":          req.City,
			"areaSqFt":      req.AreaSqFt,
			"bedrooms":      req.Bedrooms,
			"bathrooms":     req.Bathrooms,
			"amenities":     req.Amenities,
			"furnished":     req.Furnished,
			"availableFrom": req.AvailableFrom,
			"listedBy":      req.ListedBy,
			"tags":          req.Tags,
			"colorTheme":    req.ColorTheme,
			"rating":        req.Rating,
			"isVerified":    req.IsVerified,
			"listingType":   req.ListingType,
			"location":      req.Location,
			"area":          req.Area,
			"features":      req.Features,
			"status":        req.Status,
			"description":   req.Description,
			"updatedAt":     time.Now(),
		},
	}

	_, err = db.Properties.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating property"})
		return
	}

	// Clear caches
	deleteAllPropertyListCaches()
	cache.Delete("property:" + id)

	c.JSON(http.StatusOK, gin.H{"message": "Property updated successfully"})
}

func DeleteProperty(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	userID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))

	// Check if property exists and belongs to user
	var property models.Property
	err = db.Properties.FindOne(context.Background(), bson.M{
		"_id":       objectID,
		"createdBy": userID,
	}).Decode(&property)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found or unauthorized"})
		return
	}

	_, err = db.Properties.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting property"})
		return
	}

	// Clear caches
	deleteAllPropertyListCaches()
	cache.Delete("property:" + id)

	c.JSON(http.StatusOK, gin.H{"message": "Property deleted successfully"})
}

func deleteAllPropertyListCaches() {
	if cache.Client == nil {
		return
	}
	ctx := context.Background()
	keys, err := cache.Client.Keys(ctx, "properties:list*").Result()
	if err == nil {
		for _, key := range keys {
			cache.Client.Del(ctx, key)
		}
	}
}
