package accounts

type Storage interface {
	Exists(username string) bool
	Add(username string)
	Get(username string) *AccountObj
}

type GameCharacter struct {
	StartingZone string
	Name         string
	BaseColor    string
}

type AccountObj struct {
	Characters []*GameCharacter
}

type memoryStorage struct {
	storage map[string]*AccountObj
}

func NewStorage() Storage {
	return &memoryStorage{
		storage: make(map[string]*AccountObj),
	}
}

func (m *memoryStorage) Exists(username string) bool {
	_, ok := m.storage[username]
	return ok
}

func (m *memoryStorage) Add(username string) {
	m.storage[username] = &AccountObj{}
}

func (m *memoryStorage) Get(username string) *AccountObj {
	return m.storage[username]
}
