package handlers

import (
	"encoding/json"
	"github.com/UberionAI/todo-family/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

// Handler — основная структура для всех HTTP-обработчиков.
type Handler struct {
	repo *store.Store
}

// New — конструктор хэндлеров.
func New(repo *store.Store) *Handler {
	return &Handler{repo: repo}
}

// ==========================
// ======== USERS ===========
// ==========================

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TelegramID int64  `json:"telegram_id"`
		Name       string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := h.repo.CreateUser(req.TelegramID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// ==========================
// ======== GROUPS ==========
// ==========================

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string     `json:"name"`
		CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	group, err := h.repo.CreateGroup(req.Name, req.CreatedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, group)
}

func (h *Handler) RenameGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.repo.RenameGroup(id, req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "renamed"})
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeleteGroup(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.repo.ListGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, groups)
}

// ==========================
// ======== ENTRIES =========
// ==========================

func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GroupID     uuid.UUID  `json:"group_id"`
		Title       string     `json:"title"`
		Description *string    `json:"description"`
		CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	entry, err := h.repo.CreateEntry(req.GroupID, req.Title, req.Description, req.CreatedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func (h *Handler) RenameEntry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID    uuid.UUID `json:"id"`
		Title string    `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.repo.RenameEntry(req.ID, req.Title); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "renamed"})
}

func (h *Handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID uuid.UUID `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeleteEntry(req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) ListEntries(w http.ResponseWriter, r *http.Request) {
	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr == "" {
		http.Error(w, "missing group_id", http.StatusBadRequest)
		return
	}
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
		return
	}

	entries, err := h.repo.ListEntries(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

// ==========================
// ===== Helper funcs =======
// ==========================

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// форматируем даты в стиле ДД.ММ.ГГ
	type withDateFormat struct {
		Date string `json:"date,omitempty"`
	}
	_ = json.NewEncoder(w).Encode(data)
}
