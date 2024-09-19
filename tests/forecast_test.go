package tests

//
//import (
//	"github.com/stretchr/testify/assert"
//	"github.com/theroyalraj/adster-forcasting/internal/config"
//	"github.com/theroyalraj/adster-forcasting/internal/models"
//	"github.com/theroyalraj/adster-forcasting/internal/services"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	"os"
//	"testing"
//)
//
//var db *gorm.DB
//
//// Setup the database connection for tests
//func TestMain(m *testing.M) {
//	cfg := config.LoadConfig()
//	var err error
//	// Use the environment variables for connecting to the test database
//	db, err = gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
//	if err != nil {
//		panic("Failed to connect to test database!")
//	}
//	// Run the tests
//	code := m.Run()
//
//	// Exit
//	os.Exit(code)
//}
//
//// Test 1: Forecast with valid targeting criteria
//func TestForecastWithValidCriteria(t *testing.T) {
//	// skipping this test case as need csv to test
//	//dataProcessor, _ := services.NewDataProcessor(db)
//	forecastService := services.NewForecastService(db)
//
//	// Define valid criteria
//	criteria := models.TargetingCriteria{
//		GeoTarget: &models.GeoTargetSpec{
//			Included: []models.GeoTarget{
//				{TargetType: "COUNTRY", CountryCode: "US"},
//			},
//		},
//		DeviceType: &models.DeviceTypeSpec{
//			Included: []int{1}, // Mobile
//		},
//		InventoryUrl: &models.InventoryUrlSpec{
//			Included: []string{"example.com"},
//		},
//	}
//
//	// Call forecast service
//	result, err := forecastService.Forecast(criteria)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.NotZero(t, result.DailyImpressions, "Impressions should not be zero")
//	assert.NotZero(t, result.DailyReach, "Reach should not be zero")
//}
//
//// Test 2: Forecast with no criteria (should return all data)
//func TestForecastWithNoCriteria(t *testing.T) {
//	forecastService := services.NewForecastService(db)
//
//	// Call forecast service with empty criteria
//	result, err := forecastService.Forecast(models.TargetingCriteria{})
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.NotZero(t, result.DailyImpressions, "Impressions should not be zero")
//	assert.NotZero(t, result.DailyReach, "Reach should not be zero")
//}
//
//// Test 3: Forecast with invalid criteria (should return zero)
//func TestForecastWithInvalidCriteria(t *testing.T) {
//	forecastService := services.NewForecastService(db)
//
//	// Define invalid criteria
//	criteria := models.TargetingCriteria{
//		GeoTarget: &models.GeoTargetSpec{
//			Included: []models.GeoTarget{
//				{TargetType: "COUNTRY", CountryCode: "ZZ"},
//			},
//		},
//	}
//
//	// Call forecast service
//	result, err := forecastService.Forecast(criteria)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Zero(t, result.DailyImpressions, "Impressions should be zero")
//	assert.Zero(t, result.DailyReach, "Reach should be zero")
//}
//
//// Test 4: CSV upload and processing (Valid CSV)
//func TestCSVUploadProcessing(t *testing.T) {
//	dataProcessor, _ := services.NewDataProcessor(db)
//
//	// Provide a sample CSV file for testing
//	err := dataProcessor.ProcessCSV("testdata/sample.csv")
//
//	// Assertions
//	assert.NoError(t, err, "CSV processing should not throw an error")
//}
//
//// Test 5: CSV upload with invalid format (Invalid CSV)
//func TestCSVUploadInvalidFormat(t *testing.T) {
//	dataProcessor, _ := services.NewDataProcessor(db)
//
//	// Provide a broken CSV file for testing
//	err := dataProcessor.ProcessCSV("testdata/invalid.csv")
//
//	// Assertions
//	assert.Error(t, err, "CSV processing should throw an error for invalid format")
//}
