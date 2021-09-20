package selling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSellProduct(t *testing.T) {
	repository := mockStorage{}
	seller := NewService(&repository)
	assert.NoError(t, seller.SellProduct("7", 1))
	assert.Equal(t, repository.calledWithAmount, 1)
	assert.Equal(t, repository.calledWithID, "7")
}

type mockStorage struct {
	calledWithID     string
	calledWithAmount int
}

func (m *mockStorage) SellProduct(id string, amount int) error {
	m.calledWithID = id
	m.calledWithAmount = amount
	return nil
}
