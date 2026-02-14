package persistence

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// HotelRepository defines database operations for hotels.
type HotelRepository interface {
	// CreateHotel inserts a new hotel record and returns the persisted model.
	CreateHotel(hotel models.Hotel) (*models.Hotel, error)

	// GetHotelByID fetches a single hotel by its UUID primary key.
	GetHotelByID(id string) (*models.Hotel, error)

	// GetAllHotels returns every hotel, optionally filtered by location.
	GetAllHotels(location string) ([]models.Hotel, error)
}

// CreateHotel inserts a new hotel into the database.
func (r *RepositoryPg) CreateHotel(hotel models.Hotel) (*models.Hotel, error) {
	if err := r.gormDB.Create(&hotel).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to create hotel: %w", err)
	}

	return &hotel, nil
}

// GetHotelByID retrieves a hotel by its primary key.
func (r *RepositoryPg) GetHotelByID(id string) (*models.Hotel, error) {
	var hotel models.Hotel

	if err := r.gormDB.First(&hotel, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to get hotel with id %q: %w", id, err)
	}

	return &hotel, nil
}

// GetAllHotels returns all hotels ordered by rating descending.
// When location is non-empty it is applied as a case-insensitive partial filter.
func (r *RepositoryPg) GetAllHotels(location string) ([]models.Hotel, error) {
	query := r.gormDB.Model(&models.Hotel{})

	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	var hotels []models.Hotel
	if err := query.Order("rating DESC").Find(&hotels).Error; err != nil {
		return nil, fmt.Errorf("persistence: failed to list hotels: %w", err)
	}

	return hotels, nil
}