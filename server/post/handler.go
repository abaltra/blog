package post

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repository *Repository
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	p := h.Repository.Create()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	p := h.Repository.Update(3, "asdasd", false)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.Repository.DeleteByID("")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	p := h.Repository.List(0, 100)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p := h.Repository.GetBySlug("1234")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)

}
