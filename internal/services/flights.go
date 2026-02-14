package services

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// FlightService defines business operations for flights.
type FlightService interface {
	// CreateFlight validates and persists a new flight listing.
	CreateFlight(flight models.Flight) (*models.Flight, error)

	// GetFlight retrieves a single flight by ID.
	GetFlight(id string) (*models.Flight, error)

	// ListFlights returns flights, optionally filtered by origin and/or destination.
	ListFlights(origin, destination string) ([]models.Flight, error)
}

// CreateFlight validates the flight data then delegates to the repository.
func (s *TravelPlannerServiceImpl) CreateFlight(flight models.Flight) (*models.Flight, error) {
	// Business rule: arrival must be strictly after departure.
	if !flight.ArrivalTime.After(flight.DepartureTime) {
		return nil, fmt.Errorf("services: flight arrival_time must be after departure_time")
	}

	// Business rule: origin and destination must differ.
	if flight.Origin == flight.Destination {
		return nil, fmt.Errorf("services: flight origin and destination must be different")
	}

	// Business rule: available seats must be non-negative.
	if flight.SeatsAvailable < 0 {
		return nil, fmt.Errorf("services: flight seats_available must be >= 0, got %d", flight.SeatsAvailable)
	}

	created, err := s.repo.CreateFlight(flight)
	if err != nil {
		return nil, fmt.Errorf("services: create flight failed: %w", err)
	}

	return created, nil
}

// GetFlight retrieves a flight by its UUID.
func (s *TravelPlannerServiceImpl) GetFlight(id string) (*models.Flight, error) {
	if id == "" {
		return nil, fmt.Errorf("services: flight id must not be empty")
	}

	flight, err := s.repo.GetFlightByID(id)
	if err != nil {
		return nil, fmt.Errorf("services: get flight failed: %w", err)
	}

	return flight, nil
}

// ListFlights returns all flights, narrowed by origin/destination when provided.
func (s *TravelPlannerServiceImpl) ListFlights(origin, destination string) ([]models.Flight, error) {
	flights, err := s.repo.GetAllFlights(origin, destination)
	if err != nil {
		return nil, fmt.Errorf("services: list flights failed: %w", err)
	}

	return flights, nil
}