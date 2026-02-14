package services

import (
	"fmt"

	"github.com/namkatcedrickjumtock/travel-planner/internal/models"
)

// HotelService defines business operations for hotels.
type HotelService interface {
	// CreateHotel validates and persists a new hotel listing.
	CreateHotel(hotel models.Hotel) (*models.Hotel, error)

	// GetHotel retrieves a single hotel by ID.
	GetHotel(id string) (*models.Hotel, error)

	// ListHotels returns hotels, optionally filtered by location.
	ListHotels(location string) ([]models.Hotel, error)
}

// CreateHotel validates the hotel data then delegates to the repository.
func (s *TravelPlannerServiceImpl) CreateHotel(hotel models.Hotel) (*models.Hotel, error) {
	// Business rule: price per night must be a positive value.
	if hotel.PricePerNight <= 0 {
		return nil, fmt.Errorf("services: hotel price_per_night must be greater than 0")
	}

	// Business rule: rating must be between 0 and 5 when provided.
	if hotel.Rating < 0 || hotel.Rating > 5 {
		return nil, fmt.Errorf("services: hotel rating must be between 0 and 5, got %.2f", hotel.Rating)
	}

	created, err := s.repo.CreateHotel(hotel)
	if err != nil {
		return nil, fmt.Errorf("services: create hotel failed: %w", err)
	}

	return created, nil
}

// GetHotel retrieves a hotel by its UUID.
func (s *TravelPlannerServiceImpl) GetHotel(id string) (*models.Hotel, error) {
	if id == "" {
		return nil, fmt.Errorf("services: hotel id must not be empty")
	}

	hotel, err := s.repo.GetHotelByID(id)
	if err != nil {
		return nil, fmt.Errorf("services: get hotel failed: %w", err)
	}

	return hotel, nil
}

// ListHotels returns all hotels, narrowed by location when provided.
func (s *TravelPlannerServiceImpl) ListHotels(location string) ([]models.Hotel, error) {
	hotels, err := s.repo.GetAllHotels(location)
	if err != nil {
		return nil, fmt.Errorf("services: list hotels failed: %w", err)
	}

	return hotels, nil
}