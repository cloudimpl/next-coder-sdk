package resourceconsumption

type Feature struct {
	Id        string  `polycode:"id" json:"id"`
	Name      string  `json:"name"`
	Group     string  `json:"group"`
	UnitCost  float64 `json:"unitCost"`
	Total     float64 `json:"total"`
	Remaining float64 `json:"remaining"`
	Used      float64 `json:"used"`
}
