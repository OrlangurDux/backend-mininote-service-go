package routes

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"orlangur.link/services/mini.note/connectors"
	"orlangur.link/services/mini.note/controllers"
	middlewares "orlangur.link/services/mini.note/handlers"
)

// Routes -> define endpoints
func Routes() *mux.Router {
	MG := connectors.DbconnectMG()
	c := controllers.BaseController(MG)
	router := mux.NewRouter()

	apiNotAuth := router.PathPrefix("/api/v1").Subrouter()
	userNotAuth := apiNotAuth.PathPrefix("/users").Subrouter()
	userNotAuth.HandleFunc("/login", c.UserLoginEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/register", c.UserRegisterEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/forgot", c.UserForgotEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/check", c.UserCheckByEmailEndpoint).Methods("POST")

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middlewares.IsAuthorized)

	user := api.PathPrefix("/users").Subrouter()
	user.HandleFunc("/profile", c.UserProfileReadEndpoint).Methods("GET")
	user.HandleFunc("/profile", c.UserProfileUpdateEndpoint).Methods("PUT")
	user.HandleFunc("/profile", c.UserProfileDeleteEndpoint).Methods("DELETE")
	user.HandleFunc("/password", c.UserPasswordUpdateEndpoint).Methods("PUT")

	api.HandleFunc("/notes", c.NoteListEndpoint).Methods("GET")
	api.HandleFunc("/notes", c.NoteCreateEndpoint).Methods("POST")
	api.HandleFunc("/notes/search", c.NoteSearchEndpoint).Methods("GET")
	api.HandleFunc("/notes/{id}", c.NoteReadEndpoint).Methods("GET")
	api.HandleFunc("/notes/{id}", c.NoteUpdateEndpoint).Methods("PUT")
	api.HandleFunc("/notes/{id}", c.NoteDeleteEndpoint).Methods("DELETE")

	api.HandleFunc("/categories", c.CategoryListEndpoint).Methods("GET")
	api.HandleFunc("/categories", c.CategoryCreateEndpoint).Methods("POST")
	api.HandleFunc("/categories/{id}", c.CategoryReadEndpoint).Methods("GET")
	api.HandleFunc("/categories/{id}", c.CategoryUpdateEndpoint).Methods("PUT")
	api.HandleFunc("/categories/{id}", c.CategoryDeleteEndpoint).Methods("DELETE")

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	return router
}
