package models

type SickToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type AssetResponse struct {
	Batteries []Battery   `json:"batteries"`
	Positions []Position  `json:"positions"`
	Buttons   interface{} `json:"buttons"`
}

type Battery struct {
	Timestamp string  `json:"timestamp"`
	AssetId   string  `json:"assetId"`
	Level     int     `json:"level"`
	Voltage   float64 `json:"voltage"`
}

type Position struct {
	Timestamp       string `json:"timestamp"`
	AssetId         string `json:"assetId"`
	MapId           string `json:"mapId"`
	PositionDetails struct {
		X      float64     `json:"x"`
		Y      float64     `json:"y"`
		Z      float64     `json:"z"`
		Radius interface{} `json:"radius"`
	} `json:"position"`
	Accelerometer interface{} `json:"accelerometer"`
	Gyroscope     interface{} `json:"gyroscope"`
	Magnetometer  interface{} `json:"magnetometer"`
	Temperature   interface{} `json:"temperature"`
	Pressure      interface{} `json:"pressure"`
	Quaternions   interface{} `json:"quaternions"`
}
