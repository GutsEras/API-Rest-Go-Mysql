package controller

import (
	"database/sql"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Struct Controller
type usuarioController struct {
	usuarioUsecase usecase.UsuarioUsecase
}

func NewUsuarioController(usecase usecase.UsuarioUsecase) usuarioController {
	return usuarioController{
		usuarioUsecase: usecase,
	}
}

// @Summary Lista todos os usuários
// @Description Retorna todos os usuários registrados
// @Tags Usuarios
// @Produce json
// @Success 200 {array} model.Usuario
// @Failure 500 {object} model.Response
// @Router /usuarios [get]
func (u *usuarioController) GetUsuarios(ctx *gin.Context) {
	usuarios, err := u.usuarioUsecase.GetUsuarios()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, usuarios)
}

// @Summary Cria um novo usuário
// @Description Cria um novo usuário no banco de dados
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param usuario body model.Usuario true "Dados do novo usuário"
// @Success 201 {object} model.Usuario
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /usuario [post]
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

// @Summary Busca usuário por ID
// @Description Retorna os dados de um usuário pelo ID
// @Tags Usuarios
// @Produce json
// @Param usuarioId path int true "ID do usuário"
// @Success 200 {object} model.Usuario
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /usuario/{usuarioId} [get]
func (u *usuarioController) GetUsuarioById(ctx *gin.Context) {
	id := ctx.Param("usuarioId")
	if id == "" {
		response := model.Response{Message: "Id do Usuario nao pode ser nulo"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{Message: "Id do Usuario precisa ser um numero"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	usuario, err := u.usuarioUsecase.GetUsuarioById(usuarioId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if usuario == nil {
		response := model.Response{Message: "Usuario nao foi encontrado na base de dados"}
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	ctx.JSON(http.StatusOK, usuario)
}

// @Summary Atualiza usuário por ID
// @Description Atualiza os dados de um usuário existente
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param usuarioId path int true "ID do usuário"
// @Param usuario body model.Usuario true "Novos dados do usuário"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /usuario/{usuarioId} [put]
func (u *usuarioController) UpdateUsuarioById(ctx *gin.Context) {
	id := ctx.Param("usuarioId")
	if id == "" {
		response := model.Response{Message: "Id do Usuario nao pode ser nulo"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{Message: "Id do Usuario precisa ser um numero"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	var usuario model.Usuario
	if err := ctx.ShouldBindJSON(&usuario); err != nil {
		response := model.Response{Message: "Dados inválidos para o usuário"}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	err = u.usuarioUsecase.UpdateUsuarioById(usuarioId, &usuario)
	if err != nil {
		if err == sql.ErrNoRows {
			response := model.Response{Message: "Usuario não encontrado"}
			ctx.JSON(http.StatusNotFound, response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	response := model.Response{Message: "Usuario atualizado com sucesso"}
	ctx.JSON(http.StatusOK, response)
}

// @Summary Deleta (soft delete) um usuário por ID
// @Description Marca o usuário como inativo em vez de remover do banco
// @Tags Usuarios
// @Produce json
// @Param usuarioId path int true "ID do usuário"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /usuario/{usuarioId} [delete]
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
