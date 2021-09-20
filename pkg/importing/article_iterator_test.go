package importing

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleIterator(t *testing.T) {
	reader := strings.NewReader(`[
		{
			"art_id": "1",
			"name": "leg",
			"stock": "12"
		},
		{
			"art_id": "2",
			"name": "screw",
			"stock": "17"
		}
	]`)

	dec := json.NewDecoder(reader)

	token, err := dec.Token()
	assert.Nil(t, err)
	assert.Equal(t, json.Delim('['), token)

	repository := mockRepository{articles: []Article{}}
	it := NewArticleIterator(&repository)

	assert.Equal(t, 0, it.Length())

	err = it.Next(dec)
	assert.Nil(t, err)
	assert.Equal(t, 1, it.Length())

	err = it.Next(dec)
	assert.Nil(t, err)
	assert.Equal(t, 2, it.Length())

	err = it.Next(dec)
	assert.Error(t, err)

	err = it.BulkCreate()
	assert.Nil(t, err)

	assert.Equal(t, repository.articles[0].ID, "1")
	assert.Equal(t, repository.articles[0].Name, "leg")
	assert.Equal(t, repository.articles[0].Stock, "12")
	assert.Equal(t, repository.articles[1].ID, "2")
	assert.Equal(t, repository.articles[1].Name, "screw")
	assert.Equal(t, repository.articles[1].Stock, "17")

	it.Reset()
	assert.Equal(t, 0, it.Length())
}

func TestArticleIteratorOnEmptyArray(t *testing.T) {
	r1 := strings.NewReader(`[]`)
	d1 := json.NewDecoder(r1)

	t1, err := d1.Token()
	assert.Equal(t, json.Delim('['), t1)
	assert.Nil(t, err)

	repository := mockRepository{articles: []Article{}}
	it := NewArticleIterator(&repository)

	assert.Equal(t, 0, it.Length())

	err = it.Next(d1)
	assert.Error(t, err)

	err = it.BulkCreate()
	assert.Nil(t, err)
	assert.Equal(t, 0, it.Length())

	it.Reset()
	assert.Equal(t, 0, it.Length())
}

func TestArticleIteratorOnEmptyArrayObject(t *testing.T) {
	reader := strings.NewReader(`[{}]`)
	dec := json.NewDecoder(reader)

	token, err := dec.Token()
	assert.Equal(t, json.Delim('['), token)
	assert.Nil(t, err)

	repository := mockRepository{articles: []Article{}}
	it := NewArticleIterator(&repository)

	assert.Equal(t, 0, it.Length())

	err = it.Next(dec)
	assert.Nil(t, err)
	assert.Equal(t, 1, it.Length())

	err = it.BulkCreate()
	assert.Nil(t, err)
	assert.Equal(t, 1, it.Length())

	it.Reset()
	assert.Equal(t, 0, it.Length())
}

func TestArticleIteratorOnEmptyReader(t *testing.T) {
	reader := strings.NewReader("")
	dec := json.NewDecoder(reader)

	_, err := dec.Token()
	assert.Error(t, err)
	assert.Equal(t, "EOF", err.Error())

	repository := mockRepository{articles: []Article{}}
	it := NewArticleIterator(&repository)

	err = it.Next(dec)
	assert.Equal(t, "EOF", err.Error())
}

type mockRepository struct {
	articles []Article
}

func (r *mockRepository) CreateArticles(articles []Article) error {
	r.articles = append(r.articles, articles...)
	return nil
}

func (r *mockRepository) CreateProducts(products []Product) error {
	return nil
}
