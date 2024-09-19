package services

import (
	"github.com/theroyalraj/adster-forcasting/internal/models"
	"github.com/theroyalraj/adster-forcasting/internal/utils"
	"gorm.io/gorm"
	"log"
)

type ForecastService struct {
	DB *gorm.DB
}

func NewForecastService(db *gorm.DB) *ForecastService {
	return &ForecastService{
		DB: db,
	}
}

func (fs *ForecastService) Forecast(criteria models.TargetingCriteria) (models.ForecastResult, error) {

	// Create an instance of the ForecastData struct to hold the result
	var forecastData models.ForecastData

	// Start the base query from the request_logs table
	query := fs.DB.Table("request_logs").Select("COUNT(request_logs.id) AS impressions, COUNT(DISTINCT request_logs.user_id) AS reach")

	// Conditionally apply JOIN with users table if filtering on GeoTarget or DeviceType
	if criteria.GeoTarget != nil || criteria.DeviceType != nil {
		query = query.Joins("JOIN users ON request_logs.user_id = users.id")
	}

	// Conditionally apply JOIN with ad_details table if filtering on InventoryUrl
	if criteria.InventoryUrl != nil {
		query = query.Joins("JOIN ad_details ON request_logs.ad_id = ad_details.id")
	}

	// Apply GeoTarget filtering
	if criteria.GeoTarget != nil {
		if len(criteria.GeoTarget.Included) > 0 {
			var includedGeoCodes []string
			for _, geo := range criteria.GeoTarget.Included {
				includedGeoCodes = append(includedGeoCodes, geo.CountryCode)
			}
			query = query.Where("users.geo_country IN ?", includedGeoCodes)
		}
		if len(criteria.GeoTarget.Excluded) > 0 {
			var excludedGeoCodes []string
			for _, geo := range criteria.GeoTarget.Excluded {
				excludedGeoCodes = append(excludedGeoCodes, geo.CountryCode)
			}
			query = query.Where("users.geo_country NOT IN ?", excludedGeoCodes)
		}
	}

	// Apply DeviceType filtering
	if criteria.DeviceType != nil {
		if len(criteria.DeviceType.Included) > 0 {
			query = query.Where("users.device_type IN ?", criteria.DeviceType.Included)
		}
		if len(criteria.DeviceType.Excluded) > 0 {
			query = query.Where("users.device_type NOT IN ?", criteria.DeviceType.Excluded)
		}
	}

	// Apply InventoryUrl filtering
	if criteria.InventoryUrl != nil {
		if len(criteria.InventoryUrl.Included) > 0 {
			query = query.Where("ad_details.domain IN ?", criteria.InventoryUrl.Included)
		}
		if len(criteria.InventoryUrl.Excluded) > 0 {
			query = query.Where("ad_details.domain NOT IN ?", criteria.InventoryUrl.Excluded)
		}
	}

	// Log the final SQL query for debugging
	log.Println("Executing SQL Query: ", query.Statement.SQL.String())

	// Execute the query and scan the result into forecastData
	err := query.Scan(&forecastData).Error
	if err != nil {
		utils.Log.Error("Error executing forecast query:", err)
		return models.ForecastResult{}, err
	}

	// Return the forecast result based on the query result
	return models.ForecastResult{
		DailyImpressions: int(forecastData.Impressions),
		DailyReach:       int(forecastData.Reach),
	}, nil
}
