package repo

import (
	"go_rest_crud/internal/entity"
	"errors"
)

// TODO: я не знаю, куда мне это еще положить.
var (
	NotFoundErr = errors.New("not found")
)

type EquipmentStore interface {
    Add(name string, equipment entity.Equipment) error
    Get(name string) (entity.Equipment, error)
    Update(name string, equipment entity.Equipment) error
    List() (map[string]entity.Equipment, error)
    Remove(name string) error
}

// TODO: подумать какие операции нужны будут для Experiment.
type ExperimentStore interface {
    Add(name string, equipment entity.Experiment) error
    Get(name string) (entity.Experiment, error)
    Update(name string, equipment entity.Experiment) error
    List() (map[string]entity.Experiment, error)
    Remove(name string) error
}

// TODO: подумать какие операции нужны будут для ExperimentEquipment.
type ExperimentEquipmentStore interface {
    Add(name string, equipment entity.ExperimentEquipment) error
    Get(name string) (entity.ExperimentEquipment, error)
    Update(name string, equipment entity.ExperimentEquipment) error
    List() (map[string]entity.ExperimentEquipment, error)
    Remove(name string) error
}
