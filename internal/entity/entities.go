package entity

// Сущность.
type Equipment struct {
	ID            int      `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	CreationDate  string  `json:"creation_date"`      // NOTE: лучше это хранить как time.Time, но у меня не парсится потом это все нормально, так что пусть пока будет как string.
}

// TODO: В этой сущности нужно поле с датой.
// Дата добавления там и т.д.
type Experiment struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	CreationDate  string  `json:"creation_date"`      // NOTE: лучше это хранить как time.Time, но у меня не парсится потом это все нормально, так что пусть пока будет как string.
}

// Здесь мб можно хранить кол-во оборудования и дату добавления последнего,
// Но это скорее всего я буду делать в сущности Experiment.
// Это бы я предпочел хранить как 
type ExperimentEquipment struct {
	ExperimentID  int     `json:"experiment_id"`
	EquipmentID   int     `json:"equipment_id"`
}
