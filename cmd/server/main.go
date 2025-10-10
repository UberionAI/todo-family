package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/UberionAI/todo-family/internal/config"
	"github.com/UberionAI/todo-family/internal/handlers"
	"github.com/UberionAI/todo-family/internal/store"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// Логи запуска
	log.Printf("Starting app in %s mode on port %s", cfg.AppEnv, cfg.Port)
	log.Printf("Connecting to DB host=%s db=%s user=%s", cfg.PostgresHost, cfg.PostgresDB, cfg.PostgresUser)

	// Подключение к PostgreSQL
	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("connect to DB error: %v", err)
	}
	defer db.Close()

	// Репозиторий
	repo := store.New(db, cfg.MaxUsers)

	// Хэндлеры
	h := handlers.New(repo)

	// Роутер
	r := chi.NewRouter()

	// ===== Users =====
	r.Post("/users", h.CreateUser)
	r.Get("/users", h.ListUsers)

	// ===== Groups =====
	r.Post("/groups", h.CreateGroup)
	r.Patch("/groups/{id}", h.RenameGroup)
	r.Delete("/groups/{id}", h.DeleteGroup)
	r.Get("/groups", h.ListGroups)

	// ===== Entries =====
	r.Post("/entries", h.CreateEntry)
	r.Patch("/entries", h.RenameEntry)
	r.Delete("/entries", h.DeleteEntry)
	r.Get("/entries", h.ListEntries)

	// ===== Start server =====
	log.Printf("Server is running on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
