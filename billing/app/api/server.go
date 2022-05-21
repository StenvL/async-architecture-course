package api

import (
	"log"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/StenvL/async-architecture-course/billing/app/queue/producer"
	"github.com/StenvL/async-architecture-course/billing/app/repository"
	_ "github.com/StenvL/async-architecture-course/billing/docs"

	"github.com/gin-gonic/gin"
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

// Start @title Popug billing
// @version 1.0
// @description Billing service for popugs.

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl http://localhost:3000/oauth/authorize
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.manager Grants read and write access to managers information
// @scope.admin Grants read and write access to administrative information

// @BasePath /api
// @x-extension-openapi {"example": "value in a json format"}
func (s Server) SetupRoutes() {
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := s.engine.Group("/api")
	apiGroup.Use(s.authUser)

	apiGroup.GET("/account", s.GetUserAccountHandler)
	apiGroup.GET("/income", s.GetDailyIncome)
	apiGroup.POST("/pay", s.MakePayments)
}

func (s Server) Start() {
	if err := s.engine.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
