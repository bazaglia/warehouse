package importing

import (
	"encoding/json"
	"io"
)

// Iterator defines interface to iterate on JSON objects to be decoded
type Iterator interface {
	Next(dec *json.Decoder) error
	Reset()
	Length() int
	BulkCreate() error
}

func bulkImport(r io.Reader, it Iterator, skipTokens, batchSize int) error {
	dec := json.NewDecoder(r)

	for i := 0; i < skipTokens; i++ {
		if _, err := dec.Token(); err != nil {
			return err
		}
	}

	for i := 1; dec.More(); i++ {
		if err := it.Next(dec); err != nil {
			return err
		}

		if i == batchSize {
			if err := it.BulkCreate(); err != nil {
				return err
			}

			i = 1
			it.Reset()
		}
	}

	if it.Length() > 0 {
		return it.BulkCreate()
	}

	return nil
}
