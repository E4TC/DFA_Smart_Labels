package models

type Locations struct {
	ID   uint    `json:"id" gorm:"primary_key"`
	Text string  `json:"text"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Z    float32 `json:"z"`
}

type UpdateLocations struct {
	Text string  `json:"text"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Z    float32 `json:"z"`
}

type PressedLabel struct {
	LabelID string `json:"LabelID"`
}

type ActionStopFlash struct {
	Objects []string `json:"objectIds"`
}

type ActionSwitchPage struct {
	Objects  []string `json:"objectIds"`
	Duration uint     `json:"durationInMinutes"`
	Page     uint     `json:"page"`
}

type ActionFlash struct {
	Objects  []string `json:"objectIds"`
	Duration uint     `json:"durationInMinutes"`
	Patter   string   `json:"pattern"`
	Color    string   `json:"color"`
}

type ActionQuantityField struct {
	Quantity uint `json:"quantity"`
}

type ActionQuantity struct {
	Objects      string              `json:"objectId"`
	CustomFields ActionQuantityField `json:"customFields"`
}
