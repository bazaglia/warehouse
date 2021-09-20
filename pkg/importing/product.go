package importing

type ProductArticle struct {
	ID     string `json:"art_id"`
	Amount string `json:"amount_of"`
}

// Product defines properties for product to be imported
type Product struct {
	Name    string           `json:"name"`
	Article []ProductArticle `json:"contain_articles"`
}
