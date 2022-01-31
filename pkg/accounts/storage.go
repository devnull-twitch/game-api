package accounts

import "fmt"

type Storage interface {
	Exists(username string) bool
	Add(username string)
	Get(username string) *AccountObj
	AddItem(string, string, int) error
	GetItems(string, string) ([]int, error)
}

type GameCharacter struct {
	StartingZone string
	Name         string
	BaseColor    string
	Items        []int
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

func (m *memoryStorage) AddItem(username string, charName string, itemID int) error {
	acc, ok := m.storage[username]
	if !ok {
		return fmt.Errorf("unknown account name")
	}

	var char *GameCharacter
	for _, ac := range acc.Characters {
		if ac.Name == charName {
			char = ac
			break
		}
	}
	if char == nil {
		return fmt.Errorf("unknown character name")
	}

	char.Items = append(char.Items, itemID)
	return nil
}

func (m *memoryStorage) GetItems(username string, charName string) ([]int, error) {
	acc, ok := m.storage[username]
	if !ok {
		return nil, fmt.Errorf("unknown account name")
	}

	var char *GameCharacter
	for _, ac := range acc.Characters {
		if ac.Name == charName {
			char = ac
			break
		}
	}
	if char == nil {
		return nil, fmt.Errorf("unknown character name")
	}

	return char.Items, nil
}
