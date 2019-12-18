package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/yeadan/picture_backend/api/data"
	"github.com/yeadan/picture_backend/api/middlewares"
	"github.com/yeadan/picture_backend/api/routes"
)

func main() {

	//Nuevo mux.router y añadimos rutas
	router := mux.NewRouter().StrictSlash(true)
	routes.GetRoutesUsers(router)
	routes.GetRoutesComments(router)
	routes.GetRoutesLikes(router)
	routes.GetRoutesImages(router)
	routes.GetRoutesGallery(router)

	router.Use(middlewares.AuthUser)

	// Creamos Negroni y registamos middlewares
	middle := negroni.Classic()
	middle.UseHandler(router)

	//Iniciamos base de datos y .env
	data.InitDB()

	http.ListenAndServe(":8080", middle)
}
