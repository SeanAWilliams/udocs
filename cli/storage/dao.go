package storage

type Dao interface {
	Fetch(id string) ([]byte, error)
	FetchGlob(pattern string) []string
	Insert(id string, data []byte) error
	Delete(id string) error
	DeleteGlob(pattern string) error
	Index(id, title string, data []byte) error
	Query(query string) (*QueryResult, error)
	Drop() error
}
