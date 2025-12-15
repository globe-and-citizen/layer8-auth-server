package responsedto

type Response struct {
	IsSuccess bool        `json:"is_success"`
	Message   string      `json:"message"`
	Error     interface{} `json:"errors"`
	Data      interface{} `json:"data"`
}
