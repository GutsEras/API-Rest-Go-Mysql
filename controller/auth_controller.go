package controller

import (
	"go-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Usecase *usecase.AuthUsecase
}

func NewAuthController(uc *usecase.AuthUsecase) *AuthController {
	return &AuthController{Usecase: uc}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var credentials struct {
		Login string `json:"login"`
		Senha string `json:"senha"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	token, err := c.Usecase.Login(credentials.Login, credentials.Senha)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout efetuado com sucesso"})
}
