package repo

import (
	"go_rest_crud/internal/entity"
)

type MemStore struct {
	list map[string]entity.Equipment
}

func NewMemStore() *MemStore {
	list := make(map[string]entity.Equipment)
	return &MemStore{
		list,
	}
}

func (m *MemStore) Add(name string, equipment entity.Equipment) error {
	m.list[name] = equipment
	return nil
}

func (m *MemStore) Get(name string) (entity.Equipment, error) {

	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return entity.Equipment{}, NotFoundErr
}

func (m *MemStore) List() (map[string]entity.Equipment, error) {
	return m.list, nil
}

func (m *MemStore) Update(name string, equipment entity.Equipment) error {

	if _, ok := m.list[name]; ok {
		m.list[name] = equipment
		return nil
	}

	return NotFoundErr
}

func (m *MemStore) Remove(name string) error {
	delete(m.list, name)
	return nil
}
