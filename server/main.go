package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	custommiddleware "github.com/momoh-yusuf/note-app/Custom_middleware"
	"github.com/momoh-yusuf/note-app/config"
	authservice "github.com/momoh-yusuf/note-app/services/auth_service"
	noteservice "github.com/momoh-yusuf/note-app/services/note_service"
	"github.com/momoh-yusuf/note-app/utils"
)

func main() {
	// load env variables
	godotenv.Load(".env")
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// Accessing env variables
	godotenv.Load(".env")
	PORT := os.Getenv("PORT")

	apiRouter := chi.NewRouter()

	// defining auth routes
	apiRouter.Route("/auth", func(r chi.Router) {
		r.Post("/register", authservice.HandleUserRegister)
		r.Post("/login", authservice.HandleUerLogin)
	})

	// defining note routes

	apiRouter.Route("/note", func(r chi.Router) {
		r.Use(custommiddleware.AuthenticateUser)
		r.Post("/create", noteservice.HandleNoteCreation)
	})

	router.Mount("/api", apiRouter)
	router.NotFound(utils.HandleNotFound) // for handling 404
	// for starting server
	config.Db_Query()
	fmt.Println("Server Running on port number: 3000")
	http.ListenAndServe(":"+PORT, router)
}
