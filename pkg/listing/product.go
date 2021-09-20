package listing

// Product defines properties for product to be listed
type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Stock uint32 `json:"stock"`
}
