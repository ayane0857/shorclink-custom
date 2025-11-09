package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shorclick/handlers"
	"shorclick/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, proceeding with environment variables")
	}
	original_link := os.Getenv("ORIGINAL_LINK")
	// DB接続の初期化
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := "5432"
	createDatabaseIfNotExists(dbHost, dbPort, dbUser, dbPassword, dbName)
	log.Println("Connecting to database")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	log.Println("Database connection established")

	// マイグレーションの実行
	err = db.AutoMigrate(&models.ShortLink{})
	if err != nil {
		log.Fatalln("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")

	// Ginのセットアップ
	r:= gin.Default()
	r.Use(func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("Origin: %s", c.Request.Header.Get("Origin"))
		log.Printf("Headers: %v", c.Request.Header)
		c.Next()
		log.Printf("Response Status: %d", c.Writer.Status())
	})
	log.Println("Starting server on :8080")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, original_link)
	})

	api := r.Group("/api")
	api.Use(handlers.RequireAPIToken())
	{
		api.GET("/", handlers.GetShortLinks(db))
		api.GET("/:id", handlers.GetShortLink(db))
		api.POST("/", handlers.SetShortCode(), handlers.PostShortLink(db))
		api.PUT("/:id", handlers.SetShortCode(), handlers.PutShortLink(db))
		api.DELETE("/:id",  handlers.DeleteShortLink(db))
	}

	r.GET("/:id", handlers.RedirectShortLink(db))

	r.Run(":8080")
}

func createDatabaseIfNotExists(host, port, user, password, dbName string) {
	// postgresデータベースに接続（デフォルトで存在する）
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable", host, user, password, port)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to postgres database:", err)
	}

	// データベースの存在確認
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	err = db.Raw(query).Scan(&exists).Error
	if err != nil {
		log.Fatalln("Failed to check if database exists:", err)
	}

	if !exists {
		log.Printf("Database '%s' does not exist. Creating...\n", dbName)
		// データベースの作成
		createQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
		err = db.Exec(createQuery).Error
		if err != nil {
			log.Fatalln("Failed to create database:", err)
		}
		log.Printf("Database '%s' created successfully\n", dbName)
	} else {
		log.Printf("Database '%s' already exists\n", dbName)
	}

	// 接続を閉じる
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Failed to get database instance:", err)
		return
	}
	sqlDB.Close()
}
