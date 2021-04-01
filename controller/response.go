package controller

type BaseResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause,omitempty"`
}

func SuccessResonse(data interface{}) BaseResponse {
	return BaseResponse{
		Success: true,
		Data:    data,
	}
}

func InvalidRequestResponse(cause string) BaseResponse {
	return BaseResponse{
		Error: &ErrorResponse{
			Message: "invalid request",
			Cause:   cause,
		},
	}
}

func InvalidCredentialsResponse() BaseResponse {
	return BaseResponse{
		Error: &ErrorResponse{
			Message: "invalid credentials",
		},
	}
}

func InternalErrorResponse() BaseResponse {
	return BaseResponse{
		Error: &ErrorResponse{
			Message: "internal server error",
		},
	}
}
