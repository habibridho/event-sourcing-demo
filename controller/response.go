package controller

type BaseResponse struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Error   ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

func InvalidRequestResponse(cause string) BaseResponse {
	return BaseResponse{
		Error: ErrorResponse{
			Message: "invalid request",
			Cause:   cause,
		},
	}
}
