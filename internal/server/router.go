package http

import (
	"github.com/UberionAI/todo-family/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
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

	return r
}
