package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/rs/cors"
	"log"
	"net/http"
	"orlangur.link/services/mini.note/models"

	"orlangur.link/services/mini.note/controllers"
	"orlangur.link/services/mini.note/routes"

	_ "orlangur.link/services/mini.note/docs"
	middlewares "orlangur.link/services/mini.note/handlers"
)

// @title Mini Note RESTful API
// @version 0.1.2
// @description This is a backend server for service mini.note resource.

// @contact.name API Support

// @host localhost:9077
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// Version print version with start program
var Version = "development"

func init() {
	flag.StringVar(&middlewares.Env, "env", "production", "current environment")
	flag.Parse()
}

func main() {
	port := middlewares.DotEnvVariable("PORT")
	models.Version = Version
	color.Blue("Version:\t" + Version)
	color.Cyan("🌏 Server running on localhost:" + port)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	router := routes.Routes()
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	handler := c.Handler(router)
	go controllers.Tasker()

	err := http.ListenAndServe(":"+port, middlewares.LogRequest(handler))
	if err != nil {
		fmt.Println(err.Error())
	}
}
