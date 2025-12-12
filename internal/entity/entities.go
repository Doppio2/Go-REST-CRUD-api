package entity

// Сущность.
type Equipment struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
}

// TODO: В этой сущности нужно поле с датой.
// Дата добавления там и т.д.
type Experiment struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
}

// Здесь мб можно хранить кол-во оборудования и дату добавления последнего,
// Но это скорее всего я буду делать в сущности Experiment.
// Это бы я предпочел хранить как 
type ExperimentEquipment struct {
	ExperimentID  int     `json:"experiment_id"`
	EquipmentID   int     `json:"equipment_id"`
}
