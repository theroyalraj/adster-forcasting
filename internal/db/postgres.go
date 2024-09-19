package db

import (
	"github.com/theroyalraj/adster-forcasting/internal/config"
	"github.com/theroyalraj/adster-forcasting/internal/utils"
	"gorm.io/driver/postgres" // Use GORM's PostgreSQL driver
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitPostgres initializes the PostgreSQL connection using GORM
func InitPostgres(cfg *config.Config) {
	var err error

	// Use GORM's Open method with the PostgreSQL driver and the DSN (Data Source Name)
	DB, err = gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		utils.Log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	// Optional: Ping the database to test the connection
	sqlDB, err := DB.DB() // Get the generic database object sql.DB to use its features
	if err != nil {
		utils.Log.Fatal("Failed to get the database object from GORM:", err)
	}

	// Perform a ping to make sure the connection is alive
	if err := sqlDB.Ping(); err != nil {
		utils.Log.Fatal("PostgreSQL ping failed:", err)
	}

	// Log success
	utils.Log.Info("Connected to PostgreSQL with GORM")
}
