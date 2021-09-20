package importing

import "encoding/json"

type ProductIterator struct {
	products   []Product
	repository Repository
}

func (pi *ProductIterator) Next(dec *json.Decoder) error {
	var p Product
	if err := dec.Decode(&p); err != nil {
		return err
	}

	pi.products = append(pi.products, p)
	return nil
}

func (pi *ProductIterator) Reset() {
	pi.products = []Product{}
}

func (pi *ProductIterator) Length() int {
	return len(pi.products)
}

func (pi *ProductIterator) BulkCreate() error {
	return pi.repository.CreateProducts(pi.products)
}

func NewProductIterator(r Repository) Iterator {
	return &ProductIterator{
		products:   []Product{},
		repository: r,
	}
}
