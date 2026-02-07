package main

import (
	"0tak2/afterhee-server/configuration"
	"0tak2/afterhee-server/controller"
	"0tak2/afterhee-server/network"
	"0tak2/afterhee-server/repository"
	"0tak2/afterhee-server/service"
	"database/sql"
	"fmt"
	"log"

	_ "0tak2/afterhee-server/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "github.com/marcboeker/go-duckdb/v2"
	"github.com/redis/go-redis/v9"
)

// import "github.com/gofiber/fiber/v2"

func createDB(dbFileName string) *sql.DB {
	db, err := sql.Open("duckdb", dbFileName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createRDB(hostAndPort string, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     hostAndPort,
		Password: password,
		DB:       0,
	})
	return rdb
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

	redisAddress := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	rdb := createRDB(redisAddress, config.RedisPassword)
	defer rdb.Close()

	schoolRepository := repository.NewSchoolRepository(db)
	cacheRepository := repository.NewCacheRepository(rdb)
	neisNetworkRequest := network.NewNEISMealRequest()
	schoolService := service.NewSchoolService(schoolRepository, cacheRepository, neisNetworkRequest)
	schoolController := controller.NewSchoolController(schoolService)

	geminiNetworkRequest := network.NewGeminiRequest()
	geminiService := service.NewGeminiService(geminiNetworkRequest)
	suggestController := controller.NewSuggestController(geminiService)

	healthController := controller.NewHealthController()

	// App
	app := fiber.New(fiber.Config{
		ErrorHandler: controller.GlobalErrorHandler,
	})
	app.Use(logger.New())

	if config.IsDevMode {
		app.Static("/static", "./static")
	}

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// API Group
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/schools", schoolController.List)
	v1.Get("/schools/meals", schoolController.ListMeals)
	v1.Post("/suggest", suggestController.Suggest)
	v1.Get("/healthcheck", healthController.Check)

	if config.IsDevMode {
		v1.Post("/test/suggest", suggestController.SuggestTest)
	}

	log.Println("listening on :" + config.Port)
	app.Listen(":" + config.Port)
}
