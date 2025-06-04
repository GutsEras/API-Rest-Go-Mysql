package controller

import (
	"database/sql"
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TarefaController struct {
	tarefaUsecase usecase.TarefaUsecase
}

func NewTarefaController(usecase usecase.TarefaUsecase) TarefaController {
	return TarefaController{
		tarefaUsecase: usecase,
	}
}

// @Summary Lista todas as tarefas
// @Description Retorna todas as tarefas cadastradas
// @Tags Tarefas
// @Produce json
// @Success 200 {array} model.Tarefa
// @Failure 500 {object} model.Response
// @Router /tarefas [get]
func (t *TarefaController) GetTarefas(ctx *gin.Context) {
	tarefas, err := t.tarefaUsecase.GetTarefas()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, tarefas)
}

// @Summary Cria uma nova tarefa
// @Description Cria uma nova tarefa no banco de dados
// @Tags Tarefas
// @Accept json
// @Produce json
// @Param tarefa body model.Tarefa true "Dados da nova tarefa"
// @Success 201 {object} model.Tarefa
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /tarefa [post]
func (t *TarefaController) CreateTarefa(ctx *gin.Context) {
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

// @Summary Busca tarefa por ID
// @Description Retorna os dados de uma tarefa pelo ID
// @Tags Tarefas
// @Produce json
// @Param tarefaId path int true "ID da tarefa"
// @Success 200 {object} model.Tarefa
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /tarefa/{tarefaId} [get]
func (t *TarefaController) GetTarefaById(ctx *gin.Context) {
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

// @Summary Atualiza tarefa por ID
// @Description Atualiza os dados de uma tarefa existente
// @Tags Tarefas
// @Accept json
// @Produce json
// @Param tarefaId path int true "ID da tarefa"
// @Param tarefa body model.Tarefa true "Novos dados da tarefa"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /tarefa/{tarefaId} [put]
func (t *TarefaController) UpdateTarefaById(ctx *gin.Context) {
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

// @Summary Deleta (soft delete) uma tarefa por ID
// @Description Marca a tarefa como inativa em vez de removê-la do banco
// @Tags Tarefas
// @Produce json
// @Param tarefaId path int true "ID da tarefa"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /tarefa/{tarefaId} [delete]
func (t *TarefaController) SoftDeleteTarefaById(ctx *gin.Context) {
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

// @Summary Lista tarefas por usuário
// @Description Retorna todas as tarefas de um usuário específico
// @Tags Tarefas
// @Produce json
// @Param usuarioId path string true "ID do usuário"
// @Success 200 {array} model.Tarefa
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /tarefausuario/{usuarioId} [get]
func (t *TarefaController) GetTarefasByUsuarioId(ctx *gin.Context) {
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
