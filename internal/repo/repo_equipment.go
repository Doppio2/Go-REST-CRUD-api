package repo

import (
	"go_rest_crud/internal/entity"
)

type EquipmentStore interface {
    Add(name string, equipment entity.Equipment) error
    Get(name string) (entity.Equipment, error)
    Update(name string, equipment entity.Equipment) error
    List() (map[string]entity.Equipment, error)
    Remove(name string) error
}
