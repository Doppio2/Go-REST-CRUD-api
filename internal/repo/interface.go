package repo

import (
	"go_rest_crud/internal/entity"
	"errors"
)


// Интерфейс для equipment хранилища.
// Поддерживает все CRUD операции.
// ExportAllToFile метод для записи результата метода List() в файл формата .csv
type EquipmentStore interface {
    Add(equipment entity.Equipment) (int, error)
    Get(id int) (entity.Equipment, error)
    Update(id int, equipment entity.Equipment) error
    List() (map[int]entity.Equipment, error)
    Remove(id int) error
	ExportAllToFile(filePath string) error 
}

// Интерфейс для experiment хранилища.
// Поддерживает все CRUD операции.
// ExportAllToFile метод для записи результата метода List() в файл формата .csv
type ExperimentStore interface {
    Add(equipment entity.Experiment) (int, error)
    Get(id int) (entity.Experiment, error)
    Update(id int, equipment entity.Experiment) error
    List() (map[int]entity.Experiment, error)
    Remove(id int) error
	ExportAllToFile(filePath string) error 
}

// Интерфейс для хранилища связной таблицы.
// Не реализует все CRUD операции, а только
// Добавление
// Удаление
// Получаение записи об научном оборудовании
// Получения всего списка оборудования, которое учавствует в эксперименте.
type ExperimentEquipmentStore interface {
	Add(experimentId int, equipmentId int) error
    Remove(experimentId int, equipmentId int) error
	GetEquipment(experimentId int, equipmentId int) (entity.Equipment, error)
    ListEquipment(experimentId int) (map[int]entity.Equipment, error)
	ExportEquipmentToFile(experimentId int, filePath string) error
}
