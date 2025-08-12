package main

import (
	"0tak2/afterhee-server/configuration"
	"0tak2/afterhee-server/controller"
	"0tak2/afterhee-server/repository"
	"0tak2/afterhee-server/service"
	"database/sql"
	"log"

	_ "0tak2/afterhee-server/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
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

// @title AfterHee API
// @version 1.0
// @description 희그 그 이후 API
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	config := configuration.GetConfiguration()

	// Dependencies
	db := createDB(config.DBPath)
	defer db.Close()

	schoolRepository := repository.NewSchoolRepository(db)
	schoolService := service.NewSchoolService(schoolRepository)
	schoolController := controller.NewSchoolController(schoolService)

	// App
	app := fiber.New()
	app.Use(logger.New())
	app.Static("/static", "./static")
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// API Group
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/schools", schoolController.List)

	log.Println("listening on :" + config.Port)
	app.Listen(":" + config.Port)
}
