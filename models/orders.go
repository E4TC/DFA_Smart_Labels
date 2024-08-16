package models

type Order struct {
	OrderID            string                    `json:"order_id"`
	Description        string                    `json:"description"`
	Customer           string                    `json:"customer"`
	Created            int64                     `json:"created"`
	CreatedHr          string                    `json:"created_hr"`
	Positions          []OrderPosition           `json:"positions"`
	Operations         []OrderOperation          `json:"operations,omitempty"`
	OperationExecution []OrderOperationExecution `json:"operations_executions,omitempty"`
}

type OrderPosition struct {
	Guid        string `json:"guid"`
	Position    int    `json:"position"`
	Description string `json:"article"`
	Quantity    int    `json:"quantity"`
	State       string `json:"state"`
}

type OrderOperation struct {
	Guid         string `json:"guid"`
	Order        string `json:"order"`
	Position     int    `json:"position"`
	OperationID  int    `json:"operation_id"`
	Description  string `json:"description"`
	CostGroup    string `json:"cost_group"`
	MachineGroup string `json:"machine_group"`
	MachineName  string `json:"machine_name"`
	State        string `json:"state"`
}

type OrderOperationExecution struct {
	Guid                 string `json:"guid"`
	Order                string `json:"order"`
	Operation            string `json:"operation"`
	OperationDescription string `json:"operation_description"`
	CostGroup            string `json:"cost_group"`
	MachineGroup         string `json:"machine_group"`
	MachineName          string `json:"machine_name"`
	Start                string `json:"start"`
	Stop                 string `json:"stop"`
}

type OrderLabel struct {
	Label           string   `json:"label"`
	Order           string   `json:"order"`
	LabelOperations []string `json:"label_operations"`
	LabelPositions  []string `json:"label_positions"`
	Comment         string   `json:"comment"`
	Timestamp       int64    `json:"timestamp,omitempty"`
	TimestampHr     string   `json:"timestamp_hr,omitempty"`
	Status          int      `json:"status,omitempty"`
}

type UpdateLabel struct {
	Label       string   `json:"ID"`
	Order       string   `json:"bautrag"`
	Positions   string   `json:"aplan"`
	OperationID []string `json:"operationID"`
	Comment     string   `json:"text"`
}

type ActionUpdateLabelList []struct {
	Objects      string      `json:"objectId"`
	CustomFields UpdateLabel `json:"customFields"`
	Tags         []string    `json:"tags"`
}

type ActionUpdateLabel struct {
	Objects      string      `json:"objectId"`
	CustomFields UpdateLabel `json:"customFields"`
}

type OrderLabelPub struct {
	Label           string   `json:"label"`
	Order           string   `json:"order"`
	LabelOperations []string `json:"label_operations"`
	LabelPositions  []string `json:"label_positions"`
	Comment         string   `json:"comment"`
	Timestamp       int64    `json:"timestamp,omitempty"`
	TimestampHr     string   `json:"timestamp_hr,omitempty"`
}
