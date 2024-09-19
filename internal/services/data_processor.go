package services

import (
	"encoding/csv"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DataProcessor struct {
	DB        *gorm.DB
	BatchSize int
}

// Struct for the Users table
type User struct {
	ID         int    `gorm:"primaryKey"`
	UserID     string `gorm:"uniqueIndex;not null"`
	IP         string
	GeoCountry string
	GeoRegion  string
	GeoCity    string
	DeviceType string
	OS         string
	Browser    string
}

// Struct for the AdDetails table
type AdDetail struct {
	ID         int    `gorm:"primaryKey"`
	Domain     string `gorm:"not null"`
	URL        string `gorm:"not null"`
	AdPosition string
	AdSize     string
}

// Struct for the RequestLogs table
type RequestLog struct {
	ID        int       `gorm:"primaryKey"`
	Timestamp time.Time `gorm:"not null"`
	UserID    int       `gorm:"not null"`
	AdID      int       `gorm:"not null"`
}

func NewDataProcessor(db *gorm.DB) (*DataProcessor, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Get batch size from the environment, with a default value of 1000
	batchSizeStr := os.Getenv("BATCH_SIZE")
	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil || batchSize <= 0 {
		batchSize = 1000 // default to 1000 if not provided or invalid
	}

	return &DataProcessor{
		DB:        db,
		BatchSize: batchSize,
	}, nil
}

func (dp *DataProcessor) ProcessCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // variable number of fields

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Initialize a batch for bulk insertion into request_logs
	var logBatch []RequestLog
	logCount := 0

	// Assuming the first record is the header
	for i, record := range records {
		if i == 0 {
			continue
		}

		// Parse record
		timestampStr := record[0]
		ip := record[1]
		userID := record[2]
		geoCountry := record[3]
		geoRegion := record[4]
		geoCity := record[5]
		deviceType := record[6]
		os := record[7]
		browser := record[8]
		domain := record[9]
		url := record[10]
		adPosition := record[11]
		adSize := record[12]

		// Parse timestamp
		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			log.Println("Invalid timestamp:", timestampStr)
			continue
		}

		// Find or insert user
		var user User
		dp.DB.FirstOrCreate(&user, User{
			UserID:     userID,
			IP:         ip,
			GeoCountry: geoCountry,
			GeoRegion:  geoRegion,
			GeoCity:    geoCity,
			DeviceType: deviceType,
			OS:         os,
			Browser:    browser,
		})

		// Find or insert ad_details
		var adDetail AdDetail
		dp.DB.FirstOrCreate(&adDetail, AdDetail{
			Domain:     domain,
			URL:        url,
			AdPosition: adPosition,
			AdSize:     adSize,
		})

		// Add the record to the log batch for request_logs
		logBatch = append(logBatch, RequestLog{
			Timestamp: timestamp,
			UserID:    user.ID,
			AdID:      adDetail.ID,
		})
		logCount++

		// If the batch size is reached, perform the batch insert
		if logCount == dp.BatchSize {
			err = dp.batchInsertLogs(logBatch)
			if err != nil {
				log.Println("Failed to insert batch:", err)
			}

			// Reset the batch and count
			logBatch = nil
			logCount = 0
		}
	}

	// Insert any remaining records
	if len(logBatch) > 0 {
		err = dp.batchInsertLogs(logBatch)
		if err != nil {
			log.Println("Failed to insert final batch:", err)
		}
	}

	return nil
}

// batchInsertLogs performs the bulk insert of request_logs using GORM
func (dp *DataProcessor) batchInsertLogs(batch []RequestLog) error {
	return dp.DB.Transaction(func(tx *gorm.DB) error {
		// Use GORM's Create to batch insert the logs
		if err := tx.Create(&batch).Error; err != nil {
			return err
		}
		return nil
	})
}
