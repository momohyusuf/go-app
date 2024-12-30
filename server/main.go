package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/momoh-yusuf/note-app/config"
	authservice "github.com/momoh-yusuf/note-app/services/auth_service"
	"github.com/momoh-yusuf/note-app/utils"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	apiRouter := chi.NewRouter()

	// defining auth routes
	apiRouter.Route("/auth", func(r chi.Router) {
		r.Post("/register", authservice.HandleUserRegister)
		r.Post("/login", authservice.HandleUerLogin)
	})

	router.Mount("/api", apiRouter)
	router.NotFound(utils.HandleNotFound) // for handling 404
	// for starting server
	config.Db_Query()
	fmt.Println("Server Running on port number: 3000")
	http.ListenAndServe(":3000", router)
}
