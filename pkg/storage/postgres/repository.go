package postgres

import (
	"context"
	"strconv"

	"github.com/bazaglia/warehouse/pkg/importing"
	"github.com/bazaglia/warehouse/pkg/listing"
	"github.com/bazaglia/warehouse/pkg/selling"
	"github.com/jackc/pgx/v4"
)

// Storage keeps data in Postgres database
type Storage struct {
	conn *pgx.Conn
}

// GetAllProducts returns all products from the database and their stocks
func (s *Storage) GetAllProducts() ([]listing.Product, error) {
	query := `
		SELECT
			products.id,
			products.name,
			MIN(FLOOR(articles.stock / products_articles.amount)) as stock
		FROM products
		JOIN products_articles ON products.id = products_articles.product_id
		JOIN articles ON articles.id = products_articles.article_id
		GROUP BY products.id
		ORDER BY id`

	rows, err := s.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]listing.Product, 0)

	for rows.Next() {
		var id int
		var name string
		var stock uint32

		if err := rows.Scan(&id, &name, &stock); err != nil {
			return nil, err
		}

		product := listing.Product{
			ID:    strconv.Itoa(id),
			Name:  name,
			Stock: stock,
		}

		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return products, err
}

// SellProduct persists updated inventory in the database for a sold product.
// Write statement is executed within a transaction explicitly to ensure atomicity
func (s *Storage) SellProduct(id string, amount int) error {
	query := `UPDATE articles a
		SET stock = stock - pa.amount * $1
		FROM products_articles pa
		WHERE pa.product_id = $2 AND a.id=pa.article_id`

	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	// rollback is safe to call since if commit is successful, this is a no-op
	defer tx.Rollback(context.Background())

	c, err := tx.Exec(context.Background(), query, amount, id)
	if err != nil {
		return err
	}

	if c.RowsAffected() == 0 {
		return selling.ErrNotFound
	}

	return tx.Commit(context.Background())
}

// CreateArticles bulk inserts articles into the database with a single statement
func (s *Storage) CreateArticles(articles []importing.Article) error {
	columns := []string{"id", "name", "stock"}
	sql := getBulkInsertSQL("articles", columns, len(articles))
	valueArgs := make([]interface{}, 0, len(columns)*len(articles))

	for _, article := range articles {
		valueArgs = append(valueArgs, article.ID, article.Name, article.Stock)
	}

	_, err := s.conn.Exec(context.Background(), sql, valueArgs...)
	return err
}

// CreateProducts bulk inserts products into the database with a single statement
func (s *Storage) CreateProducts(products []importing.Product) error {
	columns := []string{"name"}
	valueArgs := make([]interface{}, 0, len(columns)*len(products))

	for _, p := range products {
		valueArgs = append(valueArgs, p.Name)
	}

	sql := getBulkInsertSQL("products", columns, len(products)) + " RETURNING ID"
	rows, err := s.conn.Query(context.Background(), sql, valueArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var createdIds []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		createdIds = append(createdIds, id)
	}

	if rows.Err() != nil {
		return err
	}

	columns = []string{"product_id", "article_id", "amount"}
	valueArgs = make([]interface{}, 0)

	for i, id := range createdIds {
		for _, article := range products[i].Article {
			valueArgs = append(valueArgs, id, article.ID, article.Amount)
		}
	}

	sql = getBulkInsertSQL("products_articles", columns, len(valueArgs)/len(columns))
	_, err = s.conn.Exec(context.Background(), sql, valueArgs...)

	return err
}

// NewStorage initialize a Postgres storage
func NewStorage(uri string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	return &Storage{conn}, nil
}
