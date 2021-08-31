package post

import (
	"fmt"
	"net/http"

	"github.com/abaltra/blog/server/responsehandler"
)

type Handler struct {
	Repository *Repository
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	p := Post{
		Title:    "a test",
		ID:       1,
		AuthorID: "abaltra",
	}

	p.BuildSlug()

	p, err := h.Repository.Create(p)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// p := h.Repository.Update(3, "asdasd", false)
	// TODO: FIgure out how to do updates

	fmt.Println("Update: TODO!!!")
	p := Post{
		Slug: "IM A LIEEE",
	}
	responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.Repository.DeleteByID(1)
	responsehandler.EncodeJSONResponse(w, nil, http.StatusOK, nil)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := map[string]string{
		"AuthorID": "abaltra",
	}
	p, err := h.Repository.List(0, 100, filters)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p, err := h.Repository.GetBySlug("a-test")
	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}
