package persistence

import (
	"fmt"

	"gorm.io/gorm"
)

// Repository defines all database operations for the travel planner.
// Each domain area (trips, hotels, flights, activities, bookings) is
// implemented in its own file but satisfies this single interface,
// making it straightforward to swap in a mock for unit tests.
type Repository interface {
	TripRepository
	HotelRepository
	FlightRepository
	BookingRepository
}

// RepositoryPg is the PostgreSQL implementation of Repository.
// It embeds a *gorm.DB connection that is shared across all sub-repositories.
type RepositoryPg struct {
	gormDB *gorm.DB
}

// Ensure RepositoryPg fully implements Repository at compile time.
var _ Repository = (*RepositoryPg)(nil)

// NewRepository creates a RepositoryPg and verifies the DB connection is alive.
func NewRepository(db *gorm.DB) (*RepositoryPg, error) {
	if db == nil {
		return nil, fmt.Errorf("persistence: gorm DB instance must not be nil")
	}

	return &RepositoryPg{gormDB: db}, nil
}