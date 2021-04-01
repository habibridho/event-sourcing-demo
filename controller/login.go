package controller

import (
	"github.com/dgrijalva/jwt-go"
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

	if !validCredentials(request.Email, request.Password) {
		return ctx.JSON(http.StatusUnauthorized, InvalidCredentialsResponse())
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = request.Email

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	return ctx.JSON(http.StatusOK, SuccessResonse(map[string]string{
		"token": signedToken,
	}))
}

func validCredentials(email, password string) bool {
	// TODO: check to database
	return email == "habib@email.com" && password == "password"
}
