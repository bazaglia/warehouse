CREATE TABLE products (
	id SERIAL PRIMARY KEY,
	name VARCHAR (50)
);

CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    name VARCHAR (50),
    stock numeric CHECK(stock >= 0)
);

CREATE TABLE products_articles(
    product_id INT,
    article_id INT,
    amount INT,
    PRIMARY KEY (product_id, article_id),
    CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id),
    CONSTRAINT fk_article FOREIGN KEY(article_id) REFERENCES articles(id)
);

-- INSERT INTO products(name) VALUES ('Dining Chair');

-- INSERT INTO articles(name, stock) VALUES ('Leg', 12);
-- INSERT INTO articles(name, stock) VALUES ('Screw', 17);
-- INSERT INTO articles(name, stock) VALUES ('Seat', 2);
-- INSERT INTO articles(name, stock) VALUES ('Table top', 1);

-- INSERT INTO products_articles(product_id, article_id, amount) VALUES (1, 1, 4);
-- INSERT INTO products_articles(product_id, article_id, amount) VALUES (1, 2, 8);
-- INSERT INTO products_articles(product_id, article_id, amount) VALUES (1, 3, 1);