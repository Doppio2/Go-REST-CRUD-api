package entity

// Equipment represents a single lab equipment record.
type Equipment struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CreationDate string `json:"creation_date"`
}

// Experiment represents a single experiment record.
type Experiment struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CreationDate string `json:"creation_date"`
}

// ExperimentEquipment represents a link between an experiment and equipment.
type ExperimentEquipment struct {
	ExperimentID int `json:"experiment_id"`
	EquipmentID  int `json:"equipment_id"`
}
