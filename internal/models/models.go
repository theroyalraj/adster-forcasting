package models

type TargetingCriteria struct {
	GeoTarget    *GeoTargetSpec    `json:"geo_target,omitempty"`
	DeviceType   *DeviceTypeSpec   `json:"device_type,omitempty"`
	InventoryUrl *InventoryUrlSpec `json:"inventory_url,omitempty"`
}

type GeoTargetSpec struct {
	Included []GeoTarget `json:"included,omitempty"`
	Excluded []GeoTarget `json:"excluded,omitempty"`
}

type GeoTarget struct {
	TargetType  string `json:"target_type"` // "COUNTRY", "STATE", "CITY"
	TargetId    int    `json:"target_id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code,omitempty"`
}

type DeviceTypeSpec struct {
	Included []int `json:"included,omitempty"` // 1: Mobile, 2: Tablet, 3: Desktop
	Excluded []int `json:"excluded,omitempty"`
}

type InventoryUrlSpec struct {
	Included []string `json:"included,omitempty"`
	Excluded []string `json:"excluded,omitempty"`
}

type ForecastResult struct {
	DailyImpressions int `json:"daily_impressions"`
	DailyReach       int `json:"daily_reach"`
}

type ForecastData struct {
	Impressions int64 `gorm:"column:impressions"` // Number of impressions
	Reach       int64 `gorm:"column:reach"`       // Number of unique users (reach)
}
type ForecastResponse struct {
	Forecast ForecastResult `json:"forecast"`
}
