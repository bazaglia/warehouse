package listing

// Repository provides access to the product storage
type Repository interface {
	// GetAllProducts returns all products saved in storage
	GetAllProducts() ([]Product, error)
}

// Service provides listing products operation
type Service interface {
	ListProducts() ([]Product, error)
}

type service struct {
	r Repository
}

// GetProducts returns all products
func (s *service) ListProducts() ([]Product, error) {
	return s.r.GetAllProducts()
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}
