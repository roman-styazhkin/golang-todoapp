package core_http_response

type ErrResponse struct {
	Err     string `json:"error" example:"full error message"`
	Message string `json:"message" example:"short error message"`
}
