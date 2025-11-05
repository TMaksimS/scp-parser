package domain

// @Description SCP Unit information for Creating
type CreateSCPUnit struct {
	Name        string   `json:"name"`
	Class       string   `json:"class"`
	Structure   string   `json:"structure"`
	Filial      string   `json:"filial"`
	Anomaly     string   `json:"anomaly"`
	Subject     []string `json:"subject"`
	Discription string   `json:"discription"`
	SpecialCOD  string   `json:"special_cod"`
	Property    []string `json:"property"`
	Link        string   `json:"link"`
}

type GetSCPUnit struct {
	CreateSCPUnit
	ID int
}
