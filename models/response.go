package models

type ResponseData struct {
	Status  uint        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // `omitempty` omits if data is nil
}
