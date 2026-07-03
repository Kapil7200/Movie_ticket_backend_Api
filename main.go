package main

import (
	"fmt"
	"log"

	"movie_ticket/config"
	"movie_ticket/controllers"
	"movie_ticket/database"
	"movie_ticket/middleware"

	//"movie_ticket/model"
	"movie_ticket/repository"
	"movie_ticket/routes"
	"movie_ticket/service"
	"movie_ticket/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	utils.LoadEnvFile(".env")

	if err := utils.InitLogger(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	cfg := config.LoadConfig()

	db, err := database.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logrus.Errorf("Failed to get database instance: %v", err)
			return
		}

		if err := sqlDB.Close(); err != nil {
			logrus.Errorf("Failed to close database connection: %v", err)
		} else {
			logrus.Info("Database connection closed successfully")
		}
	}()

	repo := repository.NewGORMRepository(db)
	svc := service.NewService(repo)
	jwtUtil := middleware.NewJWTUtil(cfg.JWTSecret)
	controller := controllers.NewController(svc, jwtUtil)

	router := gin.Default()
	routes.RegisterRoutes(router, controller, jwtUtil.AuthMiddleware())

	address := fmt.Sprintf(":%s", cfg.Port)

	logrus.Infof("Server started on %s", address)
	fmt.Printf("Server running on %s\n", address)
	if err := router.Run(address); err != nil {
		logrus.Errorf("Server failed: %v", err)
	}

}
