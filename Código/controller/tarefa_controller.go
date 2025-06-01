package controller

import (
	"database/sql"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type tarefaController struct {
	tarefaUsecase usecase.TarefaUsecase
}

func NewTarefaController(usecase usecase.TarefaUsecase) tarefaController {
	return tarefaController{
		tarefaUsecase: usecase,
	}
}

func (t *tarefaController) GetTarefas(ctx *gin.Context) {
	tarefas, err := t.tarefaUsecase.GetTarefas()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, tarefas)
}

func (t *tarefaController) CreateTarefa(ctx *gin.Context) {
	var tarefa model.Tarefa
	if err := ctx.BindJSON(&tarefa); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedTarefa, err := t.tarefaUsecase.CreateTarefa(tarefa)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, insertedTarefa)
}

func (t *tarefaController) GetTarefaById(ctx *gin.Context) {
	id := ctx.Param("tarefaId")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Id da Tarefa não pode ser nulo"})
		return
	}

	tarefaId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Id da Tarefa precisa ser um número"})
		return
	}

	tarefa, err := t.tarefaUsecase.GetTarefaById(tarefaId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if tarefa == nil {
		ctx.JSON(http.StatusNotFound, model.Response{Message: "Tarefa não encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, tarefa)
}

func (t *tarefaController) UpdateTarefaById(ctx *gin.Context) {
	id := ctx.Param("tarefaId")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Id da Tarefa não pode ser nulo"})
		return
	}

	tarefaId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Id da Tarefa precisa ser um número"})
		return
	}

	var tarefa model.Tarefa
	if err := ctx.ShouldBindJSON(&tarefa); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Dados inválidos para a tarefa"})
		return
	}

	err = t.tarefaUsecase.UpdateTarefaById(tarefaId, &tarefa)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, model.Response{Message: "Tarefa não encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Message: "Tarefa atualizada com sucesso"})
}

func (t *tarefaController) SoftDeleteTarefaById(ctx *gin.Context) {
	id := ctx.Param("tarefaId")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Id da Tarefa não pode ser nulo"})
		return
	}

	tarefaId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Id da Tarefa precisa ser um número"})
		return
	}

	err = t.tarefaUsecase.SoftDeleteTarefaById(tarefaId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, model.Response{Message: "Tarefa não encontrada ou já deletada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{Message: "Tarefa deletada com sucesso"})
}

func (t *tarefaController) GetTarefasByUsuarioId(ctx *gin.Context) {
	usuarioId := ctx.Param("usuarioId")
	if usuarioId == "" {
		ctx.JSON(http.StatusBadRequest, model.Response{Message: "Id do Usuário não pode ser nulo"})
		return
	}

	tarefas, err := t.tarefaUsecase.GetTarefasByUsuarioId(usuarioId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tarefas)
}
