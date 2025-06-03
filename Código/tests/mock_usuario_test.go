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

func setupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	usuarioRepository := repository.NewUsuarioRepository(db)
	usuarioUsecase := usecase.NewUsuarioUseCase(usuarioRepository)
	usuarioController := controller.NewUsuarioController(usuarioUsecase)

	router.GET("/usuarios", usuarioController.GetUsuarios)
	router.POST("/usuario", usuarioController.CreateUsuario)
	router.GET("/usuario/:usuarioId", usuarioController.GetUsuarioById)
	router.PUT("/usuario/:usuarioId", usuarioController.UpdateUsuarioById)
	router.DELETE("/usuario/:usuarioId", usuarioController.SoftDeleteUsuarioById)

	return router
}

func TestUsuarioEndpoints(t *testing.T) {
	db, mock := ConnectMockDB()
	router := setupRouter(db)

	// 1. Teste POST /usuario
	t.Run("CreateUsuario", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO usuario").
			WithArgs("Teste User", "testeuser", "123456").
			WillReturnResult(sqlmock.NewResult(1, 1))

		usuario := model.Usuario{
			Nome:  "Teste User",
			Login: "testeuser",
			Senha: "123456",
		}
		body, _ := json.Marshal(usuario)

		req, _ := http.NewRequest("POST", "/usuario", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		fmt.Println("✔️ CreateUsuario OK")
	})

	// 2. Teste GET /usuarios
	t.Run("GetUsuarios", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, nome, login, senha FROM usuario").
			WillReturnRows(sqlmock.NewRows([]string{"id", "nome", "login", "senha"}).
				AddRow(1, "João", "joao", "senha123"))

		req, _ := http.NewRequest("GET", "/usuarios", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		fmt.Println("✔️ GetUsuarios OK")
	})

	// 3. Teste GET /usuario/1
	t.Run("GetUsuarioById", func(t *testing.T) {
		mock.ExpectPrepare("SELECT id, nome, login, senha FROM usuario WHERE id = ?").
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "nome", "login", "senha"}).
				AddRow(1, "João", "joao", "senha123"))

		req, _ := http.NewRequest("GET", "/usuario/1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		fmt.Println("✔️ GetUsuarioById OK")
	})

	// 4. Teste PUT /usuario/1
	t.Run("UpdateUsuarioById", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta("UPDATE usuario SET nome = ?, login = ?, senha = ? WHERE id = ?")).
			ExpectExec().
			WithArgs("User Atualizado", "usuarioatualizado", "novaSenha123", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		update := model.Usuario{
			Nome:  "User Atualizado",
			Login: "usuarioatualizado",
			Senha: "novaSenha123",
		}
		body, _ := json.Marshal(update)

		req, _ := http.NewRequest("PUT", "/usuario/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		fmt.Println("✔️ UpdateUsuarioById OK")
	})

	// 5. Teste DELETE /usuario/1
	t.Run("SoftDeleteUsuarioById", func(t *testing.T) {
		mock.ExpectPrepare(regexp.QuoteMeta("UPDATE usuario SET ativo = 'I' WHERE id = ? AND ativo = 'A'")).
			ExpectExec().
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req, _ := http.NewRequest("DELETE", "/usuario/1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		fmt.Println("✔️ SoftDeleteUsuarioById OK")
	})

	// Verifica se todas as expectativas foram atendidas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("❌ Nem todas as expectativas foram cumpridas: %v", err)
	}
}
