package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/yourusername/todo-family/internal/config"
	"github.com/yourusername/todo-family/internal/handlers"
	"github.com/yourusername/todo-family/internal/store"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	//safe logs
	log.Printf("Starting app in %s mode on port %s", cfg.AppEnv, cfg.Port)
	log.Printf("Connecting to DB host=%s db=%s user=%s", cfg.PostgresHost, cfg.PostgresDB, cfg.PostgresUser)

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("connect to DB error: %v", err)
	}
	defer db.Close()

	repo := store.New(db, cfg.MaxUsers) //User limit
	h := handlers.New(repo)

	r := chi.NewRouter()

	//API endpoints
	r.Post("/users", h.CreateUser)
	r.Get("/users", h.ListUsers)
	r.Post("/groups/{id}", h.CreateGroup)
	r.Patch("/groups/{id}", h.RenameGroup)
	r.Get("/groups", h.ListGroups)
}
