package controller

import (
	"github.com/labstack/echo"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(ctx echo.Context) error {
	var request LoginRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, InvalidRequestResponse(err.Error()))
	}

	return ctx.String(http.StatusOK, "not finished")
}
