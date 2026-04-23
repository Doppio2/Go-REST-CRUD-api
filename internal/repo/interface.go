package repo

import (
	"errors"
	"go_rest_crud/internal/entity"
)

var (
	NotFoundErr = errors.New("not found")
)

// EquipmentStore defines CRUD operations for equipment records.
type EquipmentStore interface {
	Add(equipment entity.Equipment) (int, error)
	Get(id int) (entity.Equipment, error)
	Update(id int, equipment entity.Equipment) error
	List() (map[int]entity.Equipment, error)
	Remove(id int) error
	ExportAllToFile(filePath string) error
}

// ExperimentStore defines CRUD operations for experiments.
type ExperimentStore interface {
	Add(equipment entity.Experiment) (int, error)
	Get(id int) (entity.Experiment, error)
	Update(id int, equipment entity.Experiment) error
	List() (map[int]entity.Experiment, error)
	Remove(id int) error
	ExportAllToFile(filePath string) error
}

// ExperimentEquipmentStore defines operations for the many-to-many link table.
type ExperimentEquipmentStore interface {
	Add(experimentId int, equipmentId int) error
	Remove(experimentId int, equipmentId int) error
	GetEquipment(experimentId int, equipmentId int) (entity.Equipment, error)
	ListEquipment(experimentId int) (map[int]entity.Equipment, error)
	ExportEquipmentToFile(experimentId int, filePath string) error
}
