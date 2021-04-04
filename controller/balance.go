package controller

import (
	"event-sourcing-demo/repository"
	"event-sourcing-demo/util"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func GetBalance(ctx echo.Context) error {
	userID, err := util.GetUserIDFromEchoContext(ctx)
	if err != nil {
		log.Printf("could not get user id from context: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	account, err := repository.FetchAccountByUserID(ctx.Request().Context(), userID)
	if err != nil {
		log.Printf("could not fetch account: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	if account.ID == 0 {
		return ctx.JSON(http.StatusNotFound, GenericResponse("account not found", ""))
	}

	return ctx.JSON(http.StatusOK, SuccessResonse(map[string]interface{}{
		"balance": account.Balance,
	}))
}
