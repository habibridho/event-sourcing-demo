package controller

import (
	"context"
	"event-sourcing-demo/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
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

	signedToken, err := generateSignedToken(request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	return ctx.JSON(http.StatusOK, SuccessResonse(map[string]string{
		"token": signedToken,
	}))
}

func generateSignedToken(request LoginRequest) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = request.Email

	signedToken, err := token.SignedString([]byte("secret"))
	return signedToken, err
}

func validCredentials(email, password string) bool {
	user, err := repository.FetchUserByEmail(context.Background(), email)
	if err != nil {
		return false
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}

	return true
}
