package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
	"gorm.io/gorm"
)

// createTrip handles POST /trips.
// Expects a JSON body matching models.CreateTripRequest.
func (h *handler) createTrip(ctx *gin.Context) {
	var req models.CreateTripRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	trip, err := h.svc.CreateTrip(req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, trip)
}

// getTrip handles GET /trips/:id.
func (h *handler) getTrip(ctx *gin.Context) {
	id := ctx.Param("id")

	trip, err := h.svc.GetTrip(id)
	if err != nil {
		// Distinguish between "not found" and unexpected DB errors.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "trip not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, trip)
}

// listTrips handles GET /trips.
func (h *handler) listTrips(ctx *gin.Context) {
	trips, err := h.svc.ListTrips()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, trips)
}

// updateTrip handles PUT /trips/:id.
// Accepts a partial JSON body; only provided fields are updated.
func (h *handler) updateTrip(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateTripRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	updated, err := h.svc.UpdateTrip(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "trip not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

// deleteTrip handles DELETE /trips/:id.
func (h *handler) deleteTrip(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.svc.DeleteTrip(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "trip not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// 204 No Content — successful deletion with no body.
	ctx.Status(http.StatusNoContent)
}

// searchTrips handles GET /trips/search.
// Accepts query params: destination, start_date, end_date (YYYY-MM-DD).
func (h *handler) searchTrips(ctx *gin.Context) {
	var params models.TripSearchParams

	// ShouldBindQuery populates the struct from query-string values.
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid query parameters: " + err.Error(),
		})
		return
	}

	trips, err := h.svc.SearchTrips(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, trips)
}