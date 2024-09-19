package api

import (
	"net/http"
	"time"

	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/theroyalraj/adster-forcasting/internal/models"
	"github.com/theroyalraj/adster-forcasting/internal/services"
	"github.com/theroyalraj/adster-forcasting/internal/utils"
)

type ForecastHandler struct {
	Service     *services.ForecastService
	RedisClient *redis.Client
	Ctx         context.Context
}

func NewForecastHandler(service *services.ForecastService, redisClient *redis.Client, ctx context.Context) *ForecastHandler {
	return &ForecastHandler{
		Service:     service,
		RedisClient: redisClient,
		Ctx:         ctx,
	}
}

func (fh *ForecastHandler) Forecast(c echo.Context) error {
	var criteria models.TargetingCriteria
	if err := c.Bind(&criteria); err != nil {
		utils.Log.Error("Invalid input:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid input format",
		})
	}

	// Serialize the criteria to use as a cache key
	criteriaBytes, err := json.Marshal(criteria)
	if err != nil {
		utils.Log.Error("Failed to marshal criteria:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Internal server error",
		})
	}
	cacheKey := "forecast:" + string(criteriaBytes)

	// Check cache
	cachedResult, err := fh.RedisClient.Get(fh.Ctx, cacheKey).Result()
	if err == redis.Nil {
		// Cache miss, proceed to forecast
		utils.Log.Printf("Cache miss:", cacheKey)

		forecast, err := fh.Service.Forecast(criteria)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to generate forecast",
			})
		}

		// Cache the result with an expiration time, e.g., 1 minute
		forecastBytes, _ := json.Marshal(forecast)
		fh.RedisClient.Set(fh.Ctx, cacheKey, forecastBytes, 60*time.Second)

		return c.JSON(http.StatusOK, models.ForecastResponse{
			Forecast: forecast,
		})
	} else if err != nil {
		utils.Log.Error("Redis error:", err)
		// Proceed without cache
	} else {
		// Cache hit
		utils.Log.Printf("Cache hit:", cacheKey)
		var forecast models.ForecastResult
		if err := json.Unmarshal([]byte(cachedResult), &forecast); err != nil {
			utils.Log.Error("Failed to unmarshal cached forecast:", err)
		} else {
			return c.JSON(http.StatusOK, models.ForecastResponse{
				Forecast: forecast,
			})
		}
	}

	// If cache fails, proceed without cache
	forecast, err := fh.Service.Forecast(criteria)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to generate forecast",
		})
	}

	return c.JSON(http.StatusOK, models.ForecastResponse{
		Forecast: forecast,
	})
}
