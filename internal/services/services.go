package services

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/persistence"
)

// Planner is the top-level service interface consumed by the API layer.
// It is intentionally kept free of infrastructure concerns so it can be
// mocked easily in handler tests.
type Planner interface {
	TripService
	HotelService
	FlightService
	BookingService
}

// TravelPlannerServiceImpl is the concrete implementation of Planner.
// It delegates all data access to a persistence.Repository.
type TravelPlannerServiceImpl struct {
	repo persistence.Repository
}

// Ensure TravelPlannerServiceImpl satisfies Planner at compile time.
var _ Planner = (*TravelPlannerServiceImpl)(nil)

// NewTravelPlannerService creates a new TravelPlannerServiceImpl.
// Returns an error if repo is nil to catch wiring mistakes early.
func NewTravelPlannerService(repo persistence.Repository) (*TravelPlannerServiceImpl, error) {
	if repo == nil {
		return nil, fmt.Errorf("services: repository must not be nil")
	}

	return &TravelPlannerServiceImpl{repo: repo}, nil
}