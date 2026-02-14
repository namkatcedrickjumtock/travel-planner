package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
	"gorm.io/gorm"
)

// createFlight handles POST /flights.
// Expects a JSON body matching models.Flight (without id/created_at/updated_at).
func (h *handler) createFlight(ctx *gin.Context) {
	var flight models.Flight

	if err := ctx.ShouldBindJSON(&flight); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	created, err := h.svc.CreateFlight(flight)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

// getFlight handles GET /flights/:id.
func (h *handler) getFlight(ctx *gin.Context) {
	id := ctx.Param("id")

	flight, err := h.svc.GetFlight(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "flight not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, flight)
}

// listFlights handles GET /flights.
// Accepts optional query params: origin, destination (partial, case-insensitive).
func (h *handler) listFlights(ctx *gin.Context) {
	origin := ctx.Query("origin")
	destination := ctx.Query("destination")

	flights, err := h.svc.ListFlights(origin, destination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, flights)
}