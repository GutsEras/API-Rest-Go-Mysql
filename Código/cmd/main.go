package main

import (
	"go-api/controller"
	"go-api/db"
	docs "go-api/docs"
	"go-api/repository"
	"go-api/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// @title           API de Usuários e Tarefas
	// @version         1.0
	// @description     Esta é a API para gerenciamento de usuários e tarefas.
	// @termsOfService  http://swagger.io/terms/

	// @contact.name   Suporte da API
	// @contact.url    http://www.exemplo.com/suporte
	// @contact.email  suporte@exemplo.com

	// @license.name  Licença Apache 2.0
	// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

	// @host      localhost:8000
	// @BasePath  /

	docs.SwaggerInfo.BasePath = "/"
	server := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	// camada de repository
	UsuarioRepository := repository.NewUsuarioRepository(dbConnection)
	TarefaRepository := repository.NewTarefaRepository(dbConnection)

	// camada usecase
	UsuarioUseCase := usecase.NewUsuarioUseCase(UsuarioRepository)
	TarefaUseCase := usecase.NewTarefaUseCase(TarefaRepository)
	AuthUseCase := usecase.NewAuthUsecase(UsuarioRepository)

	// camada de controllers
	usuarioController := controller.NewUsuarioController(UsuarioUseCase)
	tarefaController := controller.NewTarefaController(TarefaUseCase)
	authController := controller.NewAuthController(AuthUseCase)

	auth := server.Group("/auth")

	// Rota de teste
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Rotas de usuário
	server.GET("/usuarios", usuarioController.GetUsuarios)
	server.POST("/usuario", usuarioController.CreateUsuario)
	server.GET("/usuario/:usuarioId", usuarioController.GetUsuarioById)
	server.PUT("/usuario/:usuarioId", usuarioController.UpdateUsuarioById)
	server.DELETE("/usuario/:usuarioId", usuarioController.SoftDeleteUsuarioById)

	// Rotas de tarefa
	server.GET("/tarefas", tarefaController.GetTarefas)
	server.POST("/tarefa", tarefaController.CreateTarefa)
	server.GET("/tarefa/:tarefaId", tarefaController.GetTarefaById)
	server.GET("/tarefausuario/:usuarioId", tarefaController.GetTarefasByUsuarioId)
	server.PUT("/tarefa/:tarefaId", tarefaController.UpdateTarefaById)
	server.DELETE("/tarefa/:tarefaId", tarefaController.SoftDeleteTarefaById)

	// Autenticação
	auth.POST("/login", authController.Login)
	auth.POST("/logout", authController.Logout)

	// Documentação Swagger
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Starta o servidor
	server.Run(":8000")

}
