package api

import (
	"log"

	"github.com/StenvL/async-architecture-course/analytics/app/repository"
	"github.com/gin-gonic/gin"

	_ "github.com/StenvL/async-architecture-course/analytics/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine *gin.Engine
	repo   repository.Repository
}

func New(repo repository.Repository) Server {
	return Server{
		engine: gin.Default(),
		repo:   repo,
	}
}

// Start @title Popug analytics
// @version 1.0
// @description Analytics service for popugs.

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl http://localhost:3000/oauth/authorize
// @scope.manager Grants manager access
// @scope.admin Grants read and write access to administrative information

// @BasePath /api
// @x-extension-openapi {"example": "value in a json format"}
func (s Server) SetupRoutes() {
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := s.engine.Group("/api")
	apiGroup.Use(s.authUser)

	apiGroup.GET("/income", s.GetDailyIncomeHandler)
	apiGroup.GET("/expensive", s.GetMostExpensiveTaskHandler)
}

func (s Server) Start() {
	if err := s.engine.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
