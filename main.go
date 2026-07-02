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

	// if err := db.AutoMigrate(&model.User{}, &model.TicketMaster{}, &model.Booking{}); err != nil {
	// 	log.Fatalf("database migration failed: %v", err)
	// }

	repo := repository.NewGORMRepository(db)
	svc := service.NewService(repo)
	jwtUtil := middleware.NewJWTUtil(cfg.JWTSecret)
	controller := controllers.NewController(svc, jwtUtil)

	router := gin.Default()
	routes.RegisterRoutes(router, controller, jwtUtil.AuthMiddleware())

	address := fmt.Sprintf(":%s", cfg.Port)

	utils.Logger.Printf("Server started on %s", address)
	fmt.Printf("Server running on %s\n", address)
	if err := router.Run(address); err != nil {
		utils.Logger.Fatalf("server failed: %v", err)
		log.Fatalf("server failed: %v", err)
	}
}
