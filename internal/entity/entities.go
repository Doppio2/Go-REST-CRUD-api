package entity

// Сущность.
type Equipment struct {
	ID            int     `json:"id"`     // TODO: узнать, нужно ли это вообще.
	Name          string  `json:"name"`
	Description   string  `json:"field"`
}

type Experiment struct {
	ID            int     `json:"id"`     // TODO: узнать, нужно ли это вообще.
	Name          string  `json:"name"`
	Description   string  `json:"field"`            
}

type ExperimentEquipment struct {
	ExperimentID  int     `json:"experiment_id"`
	EquipmentID   int     `json:"equipment_id"`
}
