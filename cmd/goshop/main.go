package main

import (
	"fmt"
	"log"

	"github.com/Alieksieiev0/goshop/internal/database"
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/Alieksieiev0/goshop/internal/transport/rest"
	"github.com/gin-gonic/gin"
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
	err = database.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)

	s := rest.NewServer(gin.New(), services.NewProductDBService(db))
	err = s.Start(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
