package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
	"gorm.io/gorm"
)

// createHotel handles POST /hotels.
// Expects a JSON body matching models.Hotel (without id/created_at/updated_at).
func (h *handler) createHotel(ctx *gin.Context) {
	var hotel models.Hotel

	if err := ctx.ShouldBindJSON(&hotel); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid request body: " + err.Error(),
		})
		return
	}

	created, err := h.svc.CreateHotel(hotel)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

// getHotel handles GET /hotels/:id.
func (h *handler) getHotel(ctx *gin.Context) {
	id := ctx.Param("id")

	hotel, err := h.svc.GetHotel(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "hotel not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, hotel)
}

// listHotels handles GET /hotels.
// Accepts optional query param: location (partial, case-insensitive match).
func (h *handler) listHotels(ctx *gin.Context) {
	location := ctx.Query("location")

	hotels, err := h.svc.ListHotels(location)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}