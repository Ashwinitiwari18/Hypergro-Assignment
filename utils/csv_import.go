package utils

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"property-listing/db"
	"property-listing/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ImportPropertiesFromCSV(path string) error {
	var reader *csv.Reader
	var file io.ReadCloser
	var err error

	if strings.HasPrefix(path, "http") {
		resp, err := http.Get(path)
		if err != nil {
			return fmt.Errorf("error downloading CSV: %v", err)
		}
		file = resp.Body
	} else {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening local CSV: %v", err)
		}
		file = f
	}
	defer file.Close()

	reader = csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV header: %v", err)
	}

	// Create header map for column indices
	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[col] = i
	}

	// Process each row
	var properties []interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV record: %v", err)
		}

		get := func(key string) string {
			idx, ok := headerMap[key]
			if !ok || idx >= len(record) {
				return ""
			}
			return record[idx]
		}

		price, _ := strconv.ParseFloat(get("price"), 64)
		bedrooms, _ := strconv.Atoi(get("bedrooms"))
		bathrooms, _ := strconv.Atoi(get("bathrooms"))
		areaSqFt, _ := strconv.ParseFloat(get("areaSqFt"), 64)
		amenities := []string{}
		if amenitiesStr := get("amenities"); amenitiesStr != "" {
			amenities = strings.Split(amenitiesStr, ",")
			for i := range amenities {
				amenities[i] = strings.TrimSpace(amenities[i])
			}
		}
		tags := []string{}
		if tagStr := get("tags"); tagStr != "" {
			tags = strings.Split(tagStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}
		rating, _ := strconv.ParseFloat(get("rating"), 64)
		isVerified := false
		if v := strings.ToLower(get("isVerified")); v == "true" || v == "1" {
			isVerified = true
		}

		property := models.Property{
			ID:            primitive.NewObjectID(),
			Title:         get("title"),
			Type:          get("type"),
			Price:         price,
			State:         get("state"),
			City:          get("city"),
			AreaSqFt:      areaSqFt,
			Bedrooms:      bedrooms,
			Bathrooms:     bathrooms,
			Amenities:     amenities,
			Furnished:     get("furnished"),
			AvailableFrom: get("availableFrom"),
			ListedBy:      get("listedBy"),
			Tags:          tags,
			ColorTheme:    get("colorTheme"),
			Rating:        rating,
			IsVerified:    isVerified,
			ListingType:   get("listingType"),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		properties = append(properties, property)
	}

	// Insert properties in batches
	batchSize := 1000
	for i := 0; i < len(properties); i += batchSize {
		end := i + batchSize
		if end > len(properties) {
			end = len(properties)
		}

		_, err := db.Properties.InsertMany(context.Background(), properties[i:end])
		if err != nil {
			return fmt.Errorf("error inserting properties batch: %v", err)
		}
	}

	return nil
}
