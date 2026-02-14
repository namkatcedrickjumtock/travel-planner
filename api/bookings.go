package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
	"gorm.io/gorm"
)

// bookItem handles POST /trips/:id/bookings.
// Creates a booking (hotel, flight, or activity) under the given trip.
func (h *handler) bookItem(ctx *gin.Context) {
	tripID := ctx.Param("id")

	var req models.CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	booking, err := h.svc.BookItem(tripID, req)
	if err != nil {
		// Surface 404 when the trip or the referenced item is not found.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, booking)
}

// getTripBookings handles GET /trips/:id/bookings.
// Returns all bookings associated with the given trip.
func (h *handler) getTripBookings(ctx *gin.Context) {
	tripID := ctx.Param("id")

	bookings, err := h.svc.GetTripBookings(tripID)
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

	ctx.JSON(http.StatusOK, bookings)
}