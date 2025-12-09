package entity

// Сущность.
// NOTE: Пока так получается, что две эти главные сущности полностью индентичны.
// Так что я не знаю имеет ли смысл иметь две разные.
type Equipment struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
}

type Experiment struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
}

type ExperimentEquipment struct {
	ExperimentID  int     `json:"experiment_id"`
	EquipmentID   int     `json:"equipment_id"`
}
