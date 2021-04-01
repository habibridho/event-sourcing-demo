package controller

import (
	"context"
	"event-sourcing-demo/repository"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"log"
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

	userID, err := getUserID(request.Email, request.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}
	if userID == 0 {
		return ctx.JSON(http.StatusUnauthorized, InvalidCredentialsResponse())
	}

	signedToken, err := generateSignedToken(request, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalErrorResponse())
	}

	return ctx.JSON(http.StatusOK, SuccessResonse(map[string]string{
		"token": signedToken,
	}))
}

func generateSignedToken(request LoginRequest, userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = request.Email
	claims["id"] = fmt.Sprintf("%d", userID)

	signedToken, err := token.SignedString([]byte("secret"))
	return signedToken, err
}

func getUserID(email, password string) (uint, error) {
	user, err := repository.FetchUserByEmail(context.Background(), email)
	if err != nil {
		log.Printf("could not fetch user by email: %s", err.Error())
		return 0, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, nil
	}

	return user.ID, nil
}
