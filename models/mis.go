package models

type MisData struct {
	Count int         `json:"count,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
}
type MisResponse struct {
	Status string  `json:"status,omitempty"`
	Msg    string  `json:"msg,omitempty"`
	Data   MisData `json:"data,omitempty"`
}

var MisResponseFailure = MisResponse{
	Status: "1",
}
var MisResponseSuccess = MisResponse{
	Status: "0",
	Msg:    "ok",
}
