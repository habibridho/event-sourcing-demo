package controller

import (
	"github.com/labstack/echo"
	"net/http"
)

func Pay(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
}
