package storage

type Storage map[string]string

func (s Storage) Add(short string, url string) {
	s[short] = url
}

func (s Storage) Get(short string) (string, bool) {
	url, ok := s[short]
	return url, ok
}

func New() Storage {
	storage := make(map[string]string)
	return storage
}
