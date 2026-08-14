package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"orlangur.link/services/mini.note/connectors"
	"orlangur.link/services/mini.note/controllers"
	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/monitoring"
)

// Routes -> define endpoints
func Routes() *mux.Router {
	middlewares.LoadDotEnv()
	MG := connectors.DbconnectMG()
	LG := monitoring.SetupLogger()
	c := controllers.BaseController(MG, LG)
	router := mux.NewRouter()

	apiNotAuth := router.PathPrefix("/api/v1").Subrouter()
	apiNotAuth.HandleFunc("/version", c.GetVersion).Methods("GET")
	userNotAuth := apiNotAuth.PathPrefix("/users").Subrouter()
	userNotAuth.HandleFunc("/login", c.UserLoginEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/register", c.UserRegisterEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/forgot", c.UserForgotEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/check", c.UserCheckByEmailEndpoint).Methods("POST")

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middlewares.IsAuthorized)

	api.HandleFunc("/send/request", c.SendRequestEndpoint).Methods("POST")
	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("/profile", c.UserProfileReadEndpoint).Methods("GET")
	users.HandleFunc("/profile", c.UserProfileUpdateEndpoint).Methods("PUT")
	users.HandleFunc("/profile", c.UserProfileDeleteEndpoint).Methods("DELETE")
	users.HandleFunc("/password", c.UserPasswordUpdateEndpoint).Methods("PUT")

	notes := api.PathPrefix("/notes").Subrouter()
	notes.HandleFunc("", c.NoteListEndpoint).Methods("GET")
	notes.HandleFunc("", c.NoteCreateEndpoint).Methods("POST")
	notes.HandleFunc("/search", c.NoteSearchEndpoint).Methods("GET")
	notes.HandleFunc("/{id}", c.NoteReadEndpoint).Methods("GET")
	notes.HandleFunc("/{id}", c.NoteUpdateEndpoint).Methods("PUT")
	notes.HandleFunc("/{id}", c.NoteDeleteEndpoint).Methods("DELETE")
	notes.HandleFunc("/favorite/{id}", c.NoteFavoriteEndpoint).Methods("PUT")

	categories := api.PathPrefix("/categories").Subrouter()
	categories.HandleFunc("", c.CategoryListEndpoint).Methods("GET")
	categories.HandleFunc("", c.CategoryCreateEndpoint).Methods("POST")
	categories.HandleFunc("/{id}", c.CategoryReadEndpoint).Methods("GET")
	categories.HandleFunc("/{id}", c.CategoryUpdateEndpoint).Methods("PUT")
	categories.HandleFunc("/{id}", c.CategoryDeleteEndpoint).Methods("DELETE")

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	fs := http.FileServer(http.Dir("./uploaded"))
	router.PathPrefix("/uploaded/").Handler(http.StripPrefix("/uploaded/", fs))

	return router
}
