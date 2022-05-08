package api

import (
	"log"

	"github.com/StenvL/async-architecture-course/app/repository"

	"github.com/StenvL/async-architecture-course/app/queue/producer"
	"github.com/gin-gonic/gin"

	_ "github.com/StenvL/async-architecture-course/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine   *gin.Engine
	repo     repository.Repository
	producer producer.Client
}

func New(repo repository.Repository, producer producer.Client) Server {
	return Server{
		engine:   gin.Default(),
		repo:     repo,
		producer: producer,
	}
}

// Start @title Popug task-tracker
// @version 1.0
// @description Task-tracker service for popugs.

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl http://localhost:3000/oauth/authorize
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @BasePath /api
// @x-extension-openapi {"example": "value in a json format"}
func (s Server) SetupRoutes() {
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := s.engine.Group("/api")
	apiGroup.Use(s.authUser)

	taskGroup := apiGroup.Group("/tasks")
	taskGroup.GET("", s.GetUserTasksListHandler)
	taskGroup.POST("", s.NewTaskHandler)
	taskGroup.POST("/resolve/:id", s.MarkTaskResolvedHandler)
	taskGroup.POST("/shuffle", s.ShuffleTasksHandler)
}

func (s Server) Start() {
	if err := s.engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
