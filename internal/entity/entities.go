package entity

// Сущность для научного оборудования. Содержит:
// Уникальный идентификатор.
// Название оборудования.
// Описание оборудования.
// Дата добавления в базу данных.
type Equipment struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	CreationDate  string  `json:"creation_date"`
}

// Сущность для научного эксперимента. Содержит:
// Уникальный идентификатор.
// Название эксперимента.
// Описание эксперимента.
// Дата добавления в базу данных.
type Experiment struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	CreationDate  string  `json:"creation_date"`      
}

// Сущность для связной таблицы. Содержит:
// ID эксперимента. Первая часть составного ключа. Внешний ключ.
// ID научного оборудования. Вторая часть составного ключа. Внешний ключ.
type ExperimentEquipment struct {
	ExperimentID  int     `json:"experiment_id"`
	EquipmentID   int     `json:"equipment_id"`
}
