package importing

import (
	"encoding/json"
)

type ArticleIterator struct {
	articles   []Article
	repository Repository
}

func (ai *ArticleIterator) Next(dec *json.Decoder) error {
	var a Article
	if err := dec.Decode(&a); err != nil {
		return err
	}

	ai.articles = append(ai.articles, a)
	return nil
}

func (ai *ArticleIterator) Reset() {
	ai.articles = []Article{}
}

func (ai *ArticleIterator) Length() int {
	return len(ai.articles)
}

func (ai *ArticleIterator) BulkCreate() error {
	return ai.repository.CreateArticles(ai.articles)
}

func NewArticleIterator(r Repository) Iterator {
	return &ArticleIterator{
		articles:   []Article{},
		repository: r,
	}
}
