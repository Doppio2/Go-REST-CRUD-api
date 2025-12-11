package repo

import (
	"go_rest_crud/internal/entity"
	"errors"
)

// TODO: я не знаю, куда мне это еще положить.
var (
	NotFoundErr = errors.New("not found")
)

// TODO: Мб стоит это объединить с EquipmentStore, а то у них разницы немного.
// Либо просто сделать ExperimentStore совместно с ExperimentEquipmentStore.
// А то у меня тогда просто два одинаковых интерфейса.
// Больше склоняюсь к варианту с общим интерфейсом с ExperimentEquipmentStore.
// TODO: во всех методах с List() я получаю все как map[int]'Entity'. 
// Если честно это выглядит каким-то избытком, но пока лучше оставлю. Не очень хочется щас это менять. Но мб стоит.
type EquipmentStore interface {
    Add(equipment entity.Equipment) error
    Get(id int) (entity.Equipment, error)
    Update(id int, equipment entity.Equipment) error
    List() (map[int]entity.Equipment, error)
    Remove(id int) error
}

type ExperimentStore interface {
    Add(equipment entity.Experiment) error
    Get(id int) (entity.Experiment, error)
    Update(id int, equipment entity.Experiment) error
    List() (map[int]entity.Experiment, error)
    Remove(id int) error
}

type ExperimentEquipmentStore interface {
	Add(experimentId int , equipmentId int) error
    Remove(experimentId int, equipmentId int) error
	GetEquipment(experimentId int, equipmentId int) (entity.Equipment, error)
    ListEquipment(experimentId int) (map[int]entity.Equipment, error)
	// NOTE: пока что это нигде не используется. И не думаю, что я придумаю, как это использовать.
	//ListExperiments(equipmentId int) (map[int]entity.Experiment, error)
}
