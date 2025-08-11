package main

import (
	"0tak2/afterhee-server/controller"
	"0tak2/afterhee-server/repository"
	"0tak2/afterhee-server/service"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/marcboeker/go-duckdb/v2"
)

// import "github.com/gofiber/fiber/v2"

func createDB(dbFileName string) *sql.DB {
	db, err := sql.Open("duckdb", dbFileName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	// ENVs
	port := getEnv("AFTERHEE_PORT", "8080")
	dbFileName := getEnv("AFTERHEE_DUCKDB_FILENAME", "database/db.duckdb")

	// Dependencies
	db := createDB(dbFileName)
	defer db.Close()

	schoolRepository := repository.NewSchoolRepository(db)
	schoolService := service.NewSchoolService(schoolRepository)
	schoolController := controller.NewSchoolController(schoolService)

	// App
	app := fiber.New()
	app.Use(logger.New())
	app.Static("/static", "./static")

	// API Group
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/schools", schoolController.List)

	log.Println("listening on :" + port)
	app.Listen(":" + port)
}

func getEnv(envKey string, fallback string) string {
	envValue := os.Getenv(envKey)

	if envValue == "" {
		return fallback
	}

	return envValue
}
