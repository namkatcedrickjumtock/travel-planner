package services

import (
	"fmt"
	"time"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// TripService defines all business operations related to trip management.
type TripService interface {
	// CreateTrip validates the request and persists a new trip.
	CreateTrip(req models.CreateTripRequest) (*models.Trip, error)

	// GetTrip retrieves a single trip by ID.
	GetTrip(id string) (*models.Trip, error)

	// ListTrips returns all trips in the system.
	ListTrips() ([]models.Trip, error)

	// UpdateTrip applies partial updates to an existing trip.
	UpdateTrip(id string, req models.UpdateTripRequest) (*models.Trip, error)

	// DeleteTrip removes a trip and all its associated bookings (via DB cascade).
	DeleteTrip(id string) error

	// SearchTrips returns trips matching the provided filter parameters.
	SearchTrips(params models.TripSearchParams) ([]models.Trip, error)
}

// CreateTrip validates the input then delegates to the repository.
func (s *TravelPlannerServiceImpl) CreateTrip(req models.CreateTripRequest) (*models.Trip, error) {
	// Business rule: end date must be after start date.
	if !req.EndDate.After(req.StartDate) {
		return nil, fmt.Errorf("services: end_date must be after start_date")
	}

	// Business rule: trips cannot be created in the past.
	if req.StartDate.Before(time.Now().UTC().Truncate(24 * time.Hour)) {
		return nil, fmt.Errorf("services: start_date cannot be in the past")
	}

	trip := models.Trip{
		UserID:      req.UserID,
		Title:       req.Title,
		Destination: req.Destination,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      models.TripStatusPlanning,
	}

	created, err := s.repo.CreateTrip(trip)
	if err != nil {
		return nil, fmt.Errorf("services: create trip failed: %w", err)
	}

	return created, nil
}

// GetTrip retrieves a trip by its UUID.
func (s *TravelPlannerServiceImpl) GetTrip(id string) (*models.Trip, error) {
	if id == "" {
		return nil, fmt.Errorf("services: trip id must not be empty")
	}

	trip, err := s.repo.GetTripByID(id)
	if err != nil {
		return nil, fmt.Errorf("services: get trip failed: %w", err)
	}

	return trip, nil
}

// ListTrips returns all trips ordered by creation time descending.
func (s *TravelPlannerServiceImpl) ListTrips() ([]models.Trip, error) {
	trips, err := s.repo.GetAllTrips()
	if err != nil {
		return nil, fmt.Errorf("services: list trips failed: %w", err)
	}

	return trips, nil
}

// UpdateTrip builds an update map from the non-nil fields in the request
// and applies it to the trip with the given ID.
func (s *TravelPlannerServiceImpl) UpdateTrip(id string, req models.UpdateTripRequest) (*models.Trip, error) {
	if id == "" {
		return nil, fmt.Errorf("services: trip id must not be empty")
	}

	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}

	if req.Destination != nil {
		updates["destination"] = *req.Destination
	}

	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}

	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		// Nothing to update — return the current record as-is.
		return s.repo.GetTripByID(id)
	}

	// Always refresh updated_at when any field changes.
	updates["updated_at"] = time.Now().UTC()

	updated, err := s.repo.UpdateTrip(id, updates)
	if err != nil {
		return nil, fmt.Errorf("services: update trip failed: %w", err)
	}

	return updated, nil
}

// DeleteTrip removes the trip with the given ID.
// Associated bookings are deleted automatically by the DB cascade constraint.
func (s *TravelPlannerServiceImpl) DeleteTrip(id string) error {
	if id == "" {
		return fmt.Errorf("services: trip id must not be empty")
	}

	if err := s.repo.DeleteTrip(id); err != nil {
		return fmt.Errorf("services: delete trip failed: %w", err)
	}

	return nil
}

// SearchTrips delegates to the repository with the provided filter params.
func (s *TravelPlannerServiceImpl) SearchTrips(params models.TripSearchParams) ([]models.Trip, error) {
	trips, err := s.repo.SearchTrips(params)
	if err != nil {
		return nil, fmt.Errorf("services: search trips failed: %w", err)
	}

	return trips, nil
}