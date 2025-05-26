package utils

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"property-listing/db"
	"property-listing/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	propertyTypes = []string{"apartment", "house", "villa", "condo", "townhouse"}
	statuses      = []string{"available", "sold", "pending"}
	features      = []string{"parking", "gym", "pool", "security", "garden", "balcony", "fireplace", "elevator"}
	cities        = []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
)

func generateRandomProperty() models.Property {
	// Generate random features (2-4 features)
	numFeatures := rand.Intn(3) + 2
	selectedFeatures := make([]string, numFeatures)
	for i := 0; i < numFeatures; i++ {
		selectedFeatures[i] = features[rand.Intn(len(features))]
	}

	// Generate random price between 100,000 and 2,000,000
	price := float64(rand.Intn(1900000) + 100000)

	// Generate random area between 500 and 5000 sq ft
	area := float64(rand.Intn(4500) + 500)

	// Generate random bedrooms (1-6)
	bedrooms := rand.Intn(5) + 1

	// Generate random bathrooms (1-4)
	bathrooms := rand.Intn(3) + 1

	// Generate random address
	streetNumber := rand.Intn(9999) + 1
	city := cities[rand.Intn(len(cities))]
	zipCode := rand.Intn(90000) + 10000

	return models.Property{
		ID:          primitive.NewObjectID(),
		Title:       generateRandomTitle(propertyTypes[rand.Intn(len(propertyTypes))]),
		Description: generateRandomDescription(),
		Price:       price,
		Location:    generateRandomAddress(streetNumber, city, zipCode),
		Bedrooms:    bedrooms,
		Bathrooms:   bathrooms,
		Area:        area,
		Type:        propertyTypes[rand.Intn(len(propertyTypes))],
		Status:      statuses[rand.Intn(len(statuses))],
		Features:    selectedFeatures,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func generateRandomTitle(propertyType string) string {
	adjectives := []string{"Luxury", "Modern", "Spacious", "Cozy", "Elegant", "Beautiful", "Stunning", "Charming"}
	locations := []string{"Downtown", "Uptown", "City Center", "Waterfront", "Hillside", "Park View", "Garden District"}

	return adjectives[rand.Intn(len(adjectives))] + " " + propertyType + " in " + locations[rand.Intn(len(locations))]
}

func generateRandomDescription() string {
	descriptions := []string{
		"Beautiful property with modern amenities and great location.",
		"Spacious living area with natural light and contemporary design.",
		"Perfect for families with excellent neighborhood amenities.",
		"Stunning views and high-end finishes throughout.",
		"Recently renovated with premium materials and appliances.",
	}
	return descriptions[rand.Intn(len(descriptions))]
}

func generateRandomAddress(streetNumber int, city string, zipCode int) string {
	streets := []string{"Main St", "Park Ave", "Broadway", "Market St", "Lake View Dr", "Ocean Blvd", "Hill St"}
	return fmt.Sprintf("%d %s, %s, %d", streetNumber, streets[rand.Intn(len(streets))], city, zipCode)
}

func GenerateAndInsertProperties(count int) error {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create a slice to hold all properties
	properties := make([]interface{}, count)

	// Generate properties
	for i := 0; i < count; i++ {
		properties[i] = generateRandomProperty()
	}

	// Insert properties in batches of 100
	batchSize := 100
	for i := 0; i < len(properties); i += batchSize {
		end := i + batchSize
		if end > len(properties) {
			end = len(properties)
		}

		_, err := db.Properties.InsertMany(context.Background(), properties[i:end])
		if err != nil {
			log.Printf("Error inserting batch %d-%d: %v", i, end, err)
			return err
		}
		log.Printf("Successfully inserted properties %d-%d", i, end)
	}

	return nil
}
