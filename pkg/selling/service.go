package selling

import (
	"errors"
)

// ErrNotFound is used when a product could not be found.
var ErrNotFound = errors.New("product not found")

// Repository provides access to the product storage
type Repository interface {
	// SellProduct sells a product and update the inventory accordingly
	SellProduct(id string, amount int) error
}

// Service provides sell product operation
type Service interface {
	SellProduct(id string, amount int) error
}

type service struct {
	r Repository
}

// GetProducts returns all products
func (s *service) SellProduct(id string, amount int) error {
	return s.r.SellProduct(id, amount)
}

// NewService creates a selling service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}
