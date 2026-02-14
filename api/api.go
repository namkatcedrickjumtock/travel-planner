package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/namkatcedrickjumtock/travel-planner/internal/services"
)

// handler holds a reference to the service layer, shared by all route files.
type handler struct {
	svc services.Planner
}

// NewAPIListener wires up the Gin router with all routes and middleware,
// then returns the engine ready to call Run() on.
func NewAPIListener(svc services.Planner) (*gin.Engine, error) {
	router := gin.Default()

	// CORS middleware — allows all origins in development.
	// Swap cors.Default() for a custom cors.Config in production.
	router.Use(cors.Default())

	h := &handler{svc: svc}

	// Health-check — useful for load balancers and container orchestrators.
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ── Trips ────────────────────────────────────────────────────────────────
	trips := router.Group("/trips")
	{
		trips.POST("", h.createTrip)
		trips.GET("", h.listTrips)
		trips.GET("/search", h.searchTrips) // must come before /:id to avoid shadowing
		trips.GET("/:id", h.getTrip)
		trips.PUT("/:id", h.updateTrip)
		trips.DELETE("/:id", h.deleteTrip)

		// Bookings nested under their parent trip.
		trips.POST("/:id/bookings", h.bookItem)
		trips.GET("/:id/bookings", h.getTripBookings)
	}

	// ── Hotels ───────────────────────────────────────────────────────────────
	hotels := router.Group("/hotels")
	{
		hotels.POST("", h.createHotel)
		hotels.GET("", h.listHotels)
		hotels.GET("/:id", h.getHotel)
	}

	// ── Flights ──────────────────────────────────────────────────────────────
	flights := router.Group("/flights")
	{
		flights.POST("", h.createFlight)
		flights.GET("", h.listFlights)
		flights.GET("/:id", h.getFlight)
	}

	return router, nil
}