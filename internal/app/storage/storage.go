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

func (s Storage) FindByValue(url string) (string, bool) {
	for k, v := range s {
		if v == url {
			return k, true
		}
	}
	return "", false
}

func (s Storage) Size() int {
	return len(s)
}
