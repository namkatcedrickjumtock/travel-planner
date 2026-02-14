package persistence

import "gorm.io/gorm"

 
// Repository persistence methods.

type Repository interface {
}

// RepositoryPg is a postgres implementation of Repository.
type RepositoryPg struct {
	gormDB *gorm.DB
}

// This line ensures that the RepositoryPg struct implements the Repository interface.
var _ Repository = &RepositoryPg{}

func NewRepository(db *gorm.DB) (*RepositoryPg, error) {
	return &RepositoryPg{
		 
		gormDB: db,
	}, nil
}
