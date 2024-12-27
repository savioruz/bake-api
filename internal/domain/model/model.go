package model

type SuccessResponse[T any] struct {
	Data     *T        `json:"data,omitempty"`
	Paginate *Paginate `json:"paginate,omitempty"`
}

type Paginate struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

type ErrorResponse struct {
	Error map[string]interface{} `json:"error"`
}

func NewSuccessResponse[T any](data *T, paginate *Paginate) *SuccessResponse[T] {
	return &SuccessResponse[T]{
		Data:     data,
		Paginate: paginate,
	}
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: map[string]interface{}{"message": err.Error()},
	}
}
