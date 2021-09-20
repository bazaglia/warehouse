package listing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListProducts(t *testing.T) {
	p1 := Product{
		ID:    "1",
		Name:  "Dining chair",
		Stock: 10,
	}

	p2 := Product{
		ID:    "2",
		Name:  "Dining table",
		Stock: 7,
	}

	repository := mockStorage{[]Product{p1, p2}}
	lister := NewService(&repository)

	products, err := lister.ListProducts()
	assert.Nil(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, p1, products[0])
	assert.Equal(t, p2, products[1])
}

type mockStorage struct {
	products []Product
}

func (m *mockStorage) GetAllProducts() ([]Product, error) {
	return m.products, nil
}
