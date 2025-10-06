package dto

type Response struct {
	ID      int    `json:"id" example:"1"`
	Message string `json:"message" example:"user added"`
}

type ErrorResponse struct {
	Error string `json:"message" example:"error happened"`
}
