package main

import (
	"go-api/controller"
	"go-api/db"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	//camada de repository
	UsuarioRepository := repository.NewUsuarioRepository(dbConnection)
	//camada usecase
	UsuarioUseCase := usecase.NewUsuarioUseCase(UsuarioRepository)
	//camada de controllers
	usuarioController := controller.NewUsuarioController(UsuarioUseCase)

	//camada de repository
	TarefaRepository := repository.NewTarefaRepository(dbConnection)
	//camada usecase
	TarefaUseCase := usecase.NewTarefaUseCase(TarefaRepository)
	//camada de controllers
	tarefaController := controller.NewTarefaController(TarefaUseCase)

	auth := server.Group("/auth")

	AuthUseCase := usecase.NewAuthUsecase(UsuarioRepository)
	authController := controller.NewAuthController(AuthUseCase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/usuarios", usuarioController.GetUsuarios)
	server.POST("/usuario", usuarioController.CreateUsuario)
	server.GET("/usuario/:usuarioId", usuarioController.GetUsuarioById)
	server.PUT("/usuario/:usuarioId", usuarioController.UpdateUsuarioById)
	server.DELETE("/usuario/:usuarioId", usuarioController.SoftDeleteUsuarioById)

	server.GET("/tarefas", tarefaController.GetTarefas)
	server.POST("/tarefa", tarefaController.CreateTarefa)
	server.GET("/tarefa/:tarefaId", tarefaController.GetTarefaById)
	server.GET("/tarefausuario/:usuarioId", tarefaController.GetTarefasByUsuarioId)
	server.PUT("/tarefa/:tarefaId", tarefaController.UpdateTarefaById)
	server.DELETE("/tarefa/:tarefaId", tarefaController.SoftDeleteTarefaById)

	auth.POST("/login", authController.Login)
	auth.POST("/logout", authController.Logout)

	server.Run(":8000")

}
