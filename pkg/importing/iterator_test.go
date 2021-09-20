package importing

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newMockReader() io.Reader {
	return strings.NewReader(`[
		{"item": "1"}, {"item": "2"}, {"item": "3"}, {"item": "4"},
		{"item": "5"}, {"item": "6"}, {"item": "7"}, {"item": "8"},
		{"item": "9"}, {"item": "10"}, {"item": "11"}, {"item": "12"}
	]`)
}

func TestIterator(t *testing.T) {
	it := mockIterator{}

	err := bulkImport(newMockReader(), &it, 1, 10)
	assert.Nil(t, err)
	assert.Len(t, it.Items, 2)
	assert.Equal(t, it.CreatedCount, 12)
	assert.Equal(t, it.ResetCount, 1)

	it = mockIterator{}

	err = bulkImport(newMockReader(), &it, 1, 50)
	assert.Nil(t, err)
	assert.Len(t, it.Items, 12)
	assert.Equal(t, it.CreatedCount, 12)
	assert.Equal(t, it.ResetCount, 0)
}

type mockItem struct {
	Item string `json:"item"`
}

type mockIterator struct {
	Items        []mockItem
	CreatedCount int
	ResetCount   int
}

func (mi *mockIterator) Next(dec *json.Decoder) error {
	var item mockItem
	if err := dec.Decode(&item); err != nil {
		return err
	}

	mi.Items = append(mi.Items, item)
	return nil
}

func (mi *mockIterator) Reset() {
	mi.Items = []mockItem{}
	mi.ResetCount++
}

func (mi *mockIterator) Length() int {
	return len(mi.Items)
}

func (mi *mockIterator) BulkCreate() error {
	mi.CreatedCount += len(mi.Items)
	return nil
}
