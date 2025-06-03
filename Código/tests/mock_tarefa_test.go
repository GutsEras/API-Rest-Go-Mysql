package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-api/controller"
	"go-api/model"
	"go-api/repository"
	"go-api/usecase"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTarefaRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	tarefaRepository := repository.NewTarefaRepository(db)
	tarefaUsecase := usecase.NewTarefaUseCase(tarefaRepository)
	tarefaController := controller.NewTarefaController(tarefaUsecase)

	router.GET("/tarefas", tarefaController.GetTarefas)
	router.POST("/tarefa", tarefaController.CreateTarefa)
	router.GET("/tarefa/:tarefaId", tarefaController.GetTarefaById)
	router.PUT("/tarefa/:tarefaId", tarefaController.UpdateTarefaById)
	router.DELETE("/tarefa/:tarefaId", tarefaController.SoftDeleteTarefaById)
	router.GET("/tarefas/usuario/:usuarioId", tarefaController.GetTarefasByUsuarioId)

	return router
}

func TestTarefaEndpoints(t *testing.T) {
	db, mock := ConnectMockDB()
	router := setupTarefaRouter(db)

	t.Run("CreateTarefa", func(t *testing.T) {
		testCreateTarefa(t, router, mock)
	})
	t.Run("GetTarefas", func(t *testing.T) {
		testGetTarefas(t, router, mock)
	})
	t.Run("GetTarefaById", func(t *testing.T) {
		testGetTarefaById(t, router, mock)
	})
	t.Run("UpdateTarefaById", func(t *testing.T) {
		testUpdateTarefaById(t, router, mock)
	})
	t.Run("SoftDeleteTarefaById", func(t *testing.T) {
		testSoftDeleteTarefaById(t, router, mock)
	})
	t.Run("GetTarefasByUsuarioId", func(t *testing.T) {
		testGetTarefasByUsuarioId(t, router, mock)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("❌ Nem todas as expectativas foram cumpridas: %v", err)
	}
}

func testCreateTarefa(t *testing.T, router *gin.Engine, mock sqlmock.Sqlmock) {
	mock.ExpectExec("INSERT INTO tarefa").
		WithArgs("Estudar Go", "Estudar interfaces", "1", "N").
		WillReturnResult(sqlmock.NewResult(1, 1))

	tarefa := model.Tarefa{
		Nome:        "Estudar Go",
		Conteudo:    "Estudar interfaces",
		UsuarioResp: "1",
		Finalizado:  "N",
	}
	body, _ := json.Marshal(tarefa)

	req, _ := http.NewRequest("POST", "/tarefa", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	fmt.Println("✔️ CreateTarefa OK")
}

func testGetTarefas(t *testing.T, router *gin.Engine, mock sqlmock.Sqlmock) {
	mock.ExpectQuery("SELECT id, nome, conteudo, usuario_responsavel, finalizado FROM tarefa").
		WillReturnRows(sqlmock.NewRows([]string{"id", "nome", "conteudo", "usuario_responsavel", "finalizado"}).
			AddRow(1, "Estudar Go", "Estudar interfaces", "1", "N"))

	req, _ := http.NewRequest("GET", "/tarefas", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	fmt.Println("✔️ GetTarefas OK")
}

func testGetTarefaById(t *testing.T, router *gin.Engine, mock sqlmock.Sqlmock) {
	mock.ExpectPrepare("SELECT id, nome, conteudo, usuario_responsavel, finalizado FROM tarefa WHERE id = ?").
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nome", "conteudo", "usuario_responsavel", "finalizado"}).
			AddRow(1, "Estudar Go", "Estudar interfaces", "1", "N"))

	req, _ := http.NewRequest("GET", "/tarefa/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	fmt.Println("✔️ GetTarefaById OK")
}

func testUpdateTarefaById(t *testing.T, router *gin.Engine, mock sqlmock.Sqlmock) {
	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE tarefa SET nome = ?, conteudo = ?, usuario_responsavel = ?, finalizado = ? WHERE id = ?")).
		ExpectExec().
		WithArgs("Go Avançado", "Estudar reflect", "1", "S", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	tarefa := model.Tarefa{
		Nome:        "Go Avançado",
		Conteudo:    "Estudar reflect",
		UsuarioResp: "1",
		Finalizado:  "S",
	}
	body, _ := json.Marshal(tarefa)

	req, _ := http.NewRequest("PUT", "/tarefa/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	fmt.Println("✔️ UpdateTarefaById OK")
}

func testSoftDeleteTarefaById(t *testing.T, router *gin.Engine, mock sqlmock.Sqlmock) {
	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE tarefa SET ativo = 'N' WHERE id = ? AND ativo = 'A'")).
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, _ := http.NewRequest("DELETE", "/tarefa/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	fmt.Println("✔️ SoftDeleteTarefaById OK")
}

func testGetTarefasByUsuarioId(t *testing.T, router *gin.Engine, mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, nome, conteudo, usuario_responsavel, finalizado FROM tarefa WHERE usuario_responsavel = ? AND ativo = 'A'")).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "nome", "conteudo", "usuario_responsavel", "finalizado"}).
			AddRow(1, "Estudar Go", "Estudar interfaces", "1", "N"))

	req, _ := http.NewRequest("GET", "/tarefas/usuario/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	fmt.Println("✔️ GetTarefasByUsuarioId OK")
}
