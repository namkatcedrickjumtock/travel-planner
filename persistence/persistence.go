package persistence

import "gorm.io/gorm"

 
// Repository persistence methods.

type Repository interface {
}

// RepositoryPg is a PostgreSQL implementation of Repository.
type RepositoryPg struct {
	gormDB *gorm.DB
}

// Ensure RepositoryPg implements the Repository interface.
var _ Repository = (*RepositoryPg)(nil)

// NewRepository creates a new RepositoryPg instance with a GORM DB connection.
func NewRepository(db *gorm.DB) (*RepositoryPg, error) {
	return &RepositoryPg{
		gormDB: db,
	}, nil
}