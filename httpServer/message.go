package httpServer

type Response struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Data    string `json:"Data"`
}

type ResponseData struct {
	Code    string      `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}
