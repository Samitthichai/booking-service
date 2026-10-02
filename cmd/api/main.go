package main

import (
	"fmt"
	"log"
	"os"

	"booking-service/internal/handler"
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"booking-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handlers struct {
	User *handler.UserHandler
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("database connected")

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}
	log.Println("schema up to date")

	// DI chain: repository -> service -> handler
	userRepo := repository.NewGormUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	h := Handlers{
		User: userHandler,
	}
	r := setupRouter(h)

	port := os.Getenv("APP_PORT")
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func connectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSLMODE"),
		)
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func setupRouter(h Handlers) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"172.20.0.0/16"})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	api.POST("/auth/register", h.User.RegisterUser)
	return r
}
