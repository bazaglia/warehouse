package importing

import "io"

// Repository provides access to the product storage
type Repository interface {
	// CreateArticles add articles to the storage
	CreateArticles([]Article) error
	// CreateProducts add products to the storage
	CreateProducts([]Product) error
}

// Service provides importing operations for articles and products
type Service interface {
	ImportArticles(r io.Reader) error
	ImportProducts(r io.Reader) error
}

type service struct {
	r Repository
}

// ImportArticles imports inventory to a database while being memory efficient
// since it uses stream to decode a streaming array of JSON objects
// https://pkg.go.dev/encoding/json#example-Decoder.Decode-Stream
func (s *service) ImportArticles(r io.Reader) error {
	return bulkImport(r, NewArticleIterator(s.r), 3, 10)
}

// ImportProducts imports products to a database while being memory efficient
// since it uses stream to decode a streaming array of JSON objects
// https://pkg.go.dev/encoding/json#example-Decoder.Decode-Stream
func (s *service) ImportProducts(r io.Reader) error {
	return bulkImport(r, NewProductIterator(s.r), 3, 10)
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}
