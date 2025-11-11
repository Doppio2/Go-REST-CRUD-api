package entity

type Equipment struct {
	ID      int     `json:"id,omitempty"`
	Name    string  `json:"name"`
	Field   string  `json:"field"`
}
