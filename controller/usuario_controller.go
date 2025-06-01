package controller

import (
	"database/sql"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type usuarioController struct {
	usuarioUsecase usecase.UsuarioUsecase
}

func NewUsuarioController(usecase usecase.UsuarioUsecase) usuarioController {
	return usuarioController{
		usuarioUsecase: usecase,
	}
}

func (u *usuarioController) GetUsuarios(ctx *gin.Context) {

	usuarios, err := u.usuarioUsecase.GetUsuarios()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, usuarios)
}

func (u *usuarioController) CreateUsuario(ctx *gin.Context) {
	var usuario model.Usuario
	if err := ctx.BindJSON(&usuario); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedUsuario, err := u.usuarioUsecase.CreateUsuario(usuario)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, insertedUsuario)
}

func (u *usuarioController) GetUsuarioById(ctx *gin.Context) {

	id := ctx.Param("usuarioId")
	if id == "" {
		response := model.Response{
			Message: "Id do Usuario nao pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Id do Usuario precisa ser um numero",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	usuario, err := u.usuarioUsecase.GetUsuarioById(usuarioId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if usuario == nil {
		response := model.Response{
			Message: "Usuario nao foi encontrado na base de dados",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, usuario)
}

func (u *usuarioController) UpdateUsuarioById(ctx *gin.Context) {
	id := ctx.Param("usuarioId")
	if id == "" {
		response := model.Response{
			Message: "Id do Usuario nao pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Id do Usuario precisa ser um numero",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var usuario model.Usuario
	if err := ctx.ShouldBindJSON(&usuario); err != nil {
		response := model.Response{
			Message: "Dados inválidos para o usuário",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = u.usuarioUsecase.UpdateUsuarioById(usuarioId, &usuario)
	if err != nil {
		if err == sql.ErrNoRows {
			response := model.Response{
				Message: "Usuario não encontrado",
			}
			ctx.JSON(http.StatusNotFound, response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response := model.Response{
		Message: "Usuario atualizado com sucesso",
	}
	ctx.JSON(http.StatusOK, response)
}

func (u *usuarioController) SoftDeleteUsuarioById(ctx *gin.Context) {
	id := ctx.Param("usuarioId")
	if id == "" {
		response := model.Response{Message: "Id do Usuario não pode ser nulo"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{Message: "Id do Usuario precisa ser um número"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = u.usuarioUsecase.SoftDeleteUsuarioById(usuarioId)
	if err != nil {
		if err == sql.ErrNoRows {
			response := model.Response{Message: "Usuário não encontrado ou já deletado"}
			ctx.JSON(http.StatusNotFound, response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := model.Response{Message: "Usuário deletado com sucesso"}
	ctx.JSON(http.StatusOK, response)
}
