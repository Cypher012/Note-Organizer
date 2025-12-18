package main

import (
	"log"
	"os"

	db "github.com/Cypher012/OrganizeNoteAPi/internal/db"
	routes "github.com/Cypher012/OrganizeNoteAPi/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	database := db.InitializeDB()

	app := fiber.New(fiber.Config{
		AppName: "Organize Note API",
	})

	routes.Setup(app, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
