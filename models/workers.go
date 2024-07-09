package models

type Workers struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Title    string  `json:"title"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	Z        float32 `json:"z"`
	Distance float32 `json:"distance"`
	URL      string  `json:"url"`
}

type UpdateWorkers struct {
	Title    string  `json:"title"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	Z        float32 `json:"z"`
	Distance float32 `json:"distance"`
	URL      string  `json:"url"`
}

type MoveWorkers struct {
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	Z        float32 `json:"z"`
	Distance float32 `json:"distance"`
}

const (
	Inside int = iota
	Outside
	Unknown
)
