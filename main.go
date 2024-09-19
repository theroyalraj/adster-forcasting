package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/theroyalraj/adster-forcasting/internal/api"
	"github.com/theroyalraj/adster-forcasting/internal/config"
	"github.com/theroyalraj/adster-forcasting/internal/db"
	"github.com/theroyalraj/adster-forcasting/internal/services"
	"github.com/theroyalraj/adster-forcasting/internal/utils"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize Logger
	utils.InitLogger()

	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Databases
	db.InitPostgres(cfg)
	db.InitRedis(cfg)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize Forecast Service
	forecastService := services.NewForecastService(db.DB)

	// Initialize API Handlers
	forecastHandler := api.NewForecastHandler(forecastService, db.RedisClient, db.Ctx)

	// Routes
	e.POST("/forecast", forecastHandler.Forecast)

	// Health Check Route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"status": "ok"})
	})

	// CSV Upload Route (Refactored to call the new function)
	e.PUT("/upload-csv", CSVUploadHandler)

	// Start Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	utils.Log.Info("Starting server on port ", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		utils.Log.Fatal("Failed to start server:", err)
	}
}

// CSVUploadHandler handles the PUT request for CSV upload
func CSVUploadHandler(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to retrieve file")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Save the uploaded file to a temporary location
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, file.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create temporary file")
	}
	defer dst.Close()

	// Copy the file content to the temporary file
	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save file")
	}

	// Process the CSV using the DataProcessor service
	dataProcessor, err := services.NewDataProcessor(db.DB)
	if err != nil {
		log.Println("Error initializing DataProcessor:", err)
		return c.String(http.StatusInternalServerError, "Error initializing DataProcessor")
	}

	err = dataProcessor.ProcessCSV(filePath)
	if err != nil {
		log.Println("Error processing CSV:", err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error processing CSV: %v", err))
	}

	// Return success response
	return c.String(http.StatusOK, "File processed successfully")
}
