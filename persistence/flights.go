package persistence

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// FlightRepository defines database operations for flights.
type FlightRepository interface {
	// CreateFlight inserts a new flight record and returns the persisted model.
	CreateFlight(flight models.Flight) (*models.Flight, error)

	// GetFlightByID fetches a single flight by its UUID primary key.
	GetFlightByID(id string) (*models.Flight, error)

	// GetAllFlights returns flights, optionally filtered by origin and/or destination.
	GetAllFlights(origin, destination string) ([]models.Flight, error)
}

// CreateFlight inserts a new flight into the database.
func (r *RepositoryPg) CreateFlight(flight models.Flight) (*models.Flight, error) {
	if err := r.gormDB.Create(&flight).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to create flight: %w", err)
	}

	return &flight, nil
}

// GetFlightByID retrieves a flight by its primary key.
func (r *RepositoryPg) GetFlightByID(id string) (*models.Flight, error) {
	var flight models.Flight

	if err := r.gormDB.First(&flight, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to get flight with id %q: %w", id, err)
	}

	return &flight, nil
}

// GetAllFlights returns all flights ordered by departure_time ascending.
// Non-empty origin / destination values are applied as case-insensitive filters.
func (r *RepositoryPg) GetAllFlights(origin, destination string) ([]models.Flight, error) {
	query := r.gormDB.Model(&models.Flight{})

	if origin != "" {
		query = query.Where("origin ILIKE ?", "%"+origin+"%")
	}

	if destination != "" {
		query = query.Where("destination ILIKE ?", "%"+destination+"%")
	}

	var flights []models.Flight
	if err := query.Order("departure_time ASC").Find(&flights).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to list flights: %w", err)
	}

	return flights, nil
}