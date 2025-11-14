package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"

	"github.com/UberionAI/todo-family/internal/config"
	"github.com/UberionAI/todo-family/internal/handlers"
	"github.com/UberionAI/todo-family/internal/server"
	"github.com/UberionAI/todo-family/internal/store"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("connect to DB error: %v", err)
	}
	defer db.Close()

	repo := store.New(db, cfg.MaxUsers)
	h := handlers.New(repo)

	router := server.NewRouter(h)

	log.Printf("Server is running on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
	fmt.Println("test")
}
