package routes

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/httprate"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"orlangur.link/services/mini.note/connectors"
	"orlangur.link/services/mini.note/controllers"
	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/monitoring"
)

func getClientIP(r *http.Request) (string, error) {
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return httprate.CanonicalizeIP(strings.TrimSpace(cfIP)), nil
	}

	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	if remoteAddr == "" {
		return "127.0.0.1", nil
	}

	if (strings.LastIndex(remoteAddr, ":")) > strings.LastIndex(remoteAddr, "]") {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return "", err
		}
		return httprate.CanonicalizeIP(ip), nil
	}

	return httprate.CanonicalizeIP(remoteAddr), nil
}

// Routes -> define endpoints
func Routes() *mux.Router {
	middlewares.LoadDotEnv()
	MG := connectors.DbconnectMG()
	RC := connectors.RedisConnect()
	LG := monitoring.SetupLogger()
	c := controllers.BaseController(MG, RC, LG)
	router := mux.NewRouter()

	apiNotAuth := router.PathPrefix("/api/v1").Subrouter()
	apiNotAuth.HandleFunc("/version", c.GetVersion).Methods("GET")
	userNotAuth := apiNotAuth.PathPrefix("/users").Subrouter()
	userNotAuth.Use(httprate.LimitBy(5, 1*time.Second, getClientIP))
	userNotAuth.HandleFunc("/login", c.UserLoginEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/register", c.UserRegisterEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/forgot", c.UserForgotEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/check", c.UserCheckByEmailEndpoint).Methods("POST")
	userNotAuth.HandleFunc("/otp", c.User2FAVerifyEndpoint).Methods("POST")

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middlewares.IsAuthorized)

	api.HandleFunc("/send/request", c.SendRequestEndpoint).Methods("POST")
	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("/profile", c.UserProfileReadEndpoint).Methods("GET")
	users.HandleFunc("/profile", c.UserProfileUpdateEndpoint).Methods("PUT")
	users.HandleFunc("/profile", c.UserProfileDeleteEndpoint).Methods("DELETE")
	users.HandleFunc("/password", c.UserPasswordUpdateEndpoint).Methods("PUT")
	users.HandleFunc("/tfa", c.User2FAEnableEndpoint).Methods("PUT")

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
