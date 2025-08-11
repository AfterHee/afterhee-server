package main

import (
	"0tak2/afterhee-server/controller"
	"0tak2/afterhee-server/repository"
	"0tak2/afterhee-server/service"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/marcboeker/go-duckdb/v2"
)

// import "github.com/gofiber/fiber/v2"

func createDB() *sql.DB {
	db, err := sql.Open("duckdb", "database/db.duckdb")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	// Dependencies
	db := createDB()
	defer db.Close()

	schoolRepository := repository.NewSchoolRepository(db)
	schoolService := service.NewSchoolService(schoolRepository)
	schoolController := controller.NewSchoolController(schoolService)

	// App
	app := fiber.New()
	app.Use(logger.New())

	// API Group
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/schools", schoolController.List)

	log.Println("listening on :3000")
	app.Listen(":3000")
}
