package main

import (
	"log"

	"github.com/Alieksieiev0/goshop/internal/database"
	"github.com/Alieksieiev0/goshop/internal/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = database.Setup(db)
	if err != nil {
		log.Fatal(err)
	}

	s := rest.NewServer(
		fiber.New(),
		db,
	)
	err = s.Start(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
