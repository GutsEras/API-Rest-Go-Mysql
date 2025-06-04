package controller

import (
	"go-api/model"
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

// @Summary Efetua login
// @Description Realiza autenticação do usuário e retorna um token JWT
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param credentials body model.LoginRequest true "Credenciais do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var credentials model.LoginRequest
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

// @Summary Efetua logout
// @Description Simula o logout do usuário
// @Tags Autenticação
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout efetuado com sucesso"})
}
