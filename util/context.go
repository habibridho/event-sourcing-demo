package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"strconv"
)

func GetUserIDFromEchoContext(ctx echo.Context) (uint, error) {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		err := errors.New("could not extract token from context")
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("could not extract claims from token")
		return 0, err
	}
	senderIDStr, ok := claims["id"].(string)
	if !ok {
		err := errors.New("could not get user id from claims")
		return 0, err
	}
	senderID, err := strconv.ParseUint(senderIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(senderID), nil
}
