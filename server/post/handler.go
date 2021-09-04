package post

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/abaltra/blog/server/responsehandler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	Repository *Repository
}

type updateRequest struct {
	Body string `json:"Body"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
		return
	}

	var post Post
	json.Unmarshal(b, &post)

	post.BuildSlug()
	post.CreatedAt = time.Now()
	post.Version = 1
	post.ID = uuid.New().String()
	post.AuthorID = "abaltra"

	p, err := h.Repository.Create(post)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, err := h.Repository.GetBySlug(vars["slug"])

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
		return
	}

	if p == nil {
		responsehandler.EncodeJSONError(w, nil, http.StatusNotFound)
		return
	}

	if p.IsPublished {
		return
	}

	p.PublishedAt = time.Now()
	p.IsPublished = true

	h.Repository.Save(*p)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, err := h.Repository.GetBySlug(vars["slug"])

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
		return
	}

	if p == nil {
		responsehandler.EncodeJSONError(w, nil, http.StatusNotFound)
		return
	}

	requestContents, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
		return
	}

	var ur updateRequest
	json.Unmarshal(requestContents, &ur)

	p.UpdatedAt = time.Now()
	p.ContentRaw = ur.Body

	if err := h.Repository.Save(*p); err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusInternalServerError)
		return
	}

	responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, err := h.Repository.GetBySlug(vars["slug"])

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadGateway)
		return
	}

	if p == nil {
		responsehandler.EncodeJSONError(w, nil, http.StatusNotFound)
		return
	}

	h.Repository.DeleteByID(p.ID)
	responsehandler.EncodeJSONResponse(w, nil, http.StatusOK, nil)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	from := vars["from"]
	var from_int int
	if from == "" {
		from_int = 0
	} else {
		var err error
		from_int, err = strconv.Atoi(from)

		if err != nil {
			responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
			return
		}
	}

	size := vars["size"]
	var size_int int
	if size == "" {
		size_int = 100
	} else {
		var err error
		size_int, err = strconv.Atoi(size)

		if err != nil {
			responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
			return
		}
	}

	if from_int < 0 || size_int > 200 {
		responsehandler.EncodeJSONError(w, fmt.Errorf("Invalid FROM %d smaller than 0 or SIZE %d larger than 100", from_int, size_int), http.StatusBadRequest)
	}

	showDrafts, _ := strconv.ParseBool(vars["showDrafts"])

	filters := make(map[string]interface{})
	filters["IsPublished"] = true

	if showDrafts {
		filters["AuthorID"] = "abaltra"
		filters["IsPublished"] = nil
	}

	p, err := h.Repository.List(from_int, size_int, filters)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	p, err := h.Repository.GetBySlug(vars["slug"])
	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}
