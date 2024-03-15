package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alieksieiev0/goshop/internal/database"
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/Alieksieiev0/goshop/internal/transport/rest"
	"github.com/donseba/go-htmx"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	htmx *htmx.HTMX
}

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

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	// initiate a new htmx handler
	h := a.htmx.NewHandler(w, r)

	// check if the request is a htmx request
	if h.IsHxRequest() {
		// do something
	}

	// check if the request is boosted
	if h.IsHxBoosted() {
		// do something
	}

	// check if the request is a history restore request
	if h.IsHxHistoryRestoreRequest() {
		// do something
	}

	// check if the request is a prompt request
	if h.RenderPartial() {
		// do something
	}

	// set the headers for the response, see docs for more options
	swap := htmx.NewSwap().Swap().ScrollBottom()
	h.ReSwapWithObject(swap)
	// write the output like you normally do.
	// check the inspector tool in the browser to see that the headers are set.
	fmt.Println("response")
	_, _ = h.Write([]byte("OK"))
}
