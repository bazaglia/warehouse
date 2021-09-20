package importing

// Article defines properties for the article to be imported
type Article struct {
	ID    string `json:"art_id"`
	Name  string `json:"name"`
	Stock string `json:"stock"`
}
