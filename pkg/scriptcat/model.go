package scriptcat

type Metadata map[string][]string

type Script struct {
	ID       string   `json:"id"`
	Code     string   `json:"code"`
	Metadata Metadata `json:"metadata"`
}

func (s *Script) StorageName() string {
	storageNames, ok := s.Metadata["storageName"]
	if !ok {
		storageNames = []string{s.ID}
	}
	return storageNames[0]
}
