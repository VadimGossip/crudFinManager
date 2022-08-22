package domain

type ErrorResponse struct {
	Message string `json:"error" example:"error message"`
}

type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}
