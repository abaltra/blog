package post

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"glog/responsehandler"

	"github.com/gorilla/mux"
)

type Handler struct {
	Repository *Repository
}

type UpdateRequest struct {
	Body string `json:"Body"`
}

// Create godoc
// @Summary      Create a new Post
// @Description  Create a new post with an auto-generated ID
// @Accept       json
// @Produce      json
// @Param        tenantID   path      int  true  "Tenant ID"
// @Param        {object} body post.CreatePostRequest true "Post to create"
// @Success      200  {object}  post.Post
// @Failure      400  {object}  responsehandler.Error
// @Failure      500  {object}  responsehandler.Error
// @Router       /v2/tenant/{tenantID}/posts [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	vars := mux.Vars(r)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
		return
	}

	var createRequest CreatePostRequest

	json.Unmarshal(b, &createRequest)

	post := NewPost("abaltra", createRequest)

	p, err := h.Repository.Create(vars["tenantID"], *post)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}

// Publish godoc
// @Summary      Publish a Post
// @Description  Publshes an existing Post, adding the proper timestamps
// @Accept       json
// @Produce      json
// @Param        tenantID   path      int  true  "Tenant ID"
// @Param        slug   path      string  true  "Unique slug of the post"
// @Success      200
// @Failure      400  {object}  responsehandler.Error
// @Failure      500  {object}  responsehandler.Error
// @Router       /v2/tenant/{tenantID}/posts/{slug}/publish [put]
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, err := h.Repository.GetBySlug(vars["tenantID"], vars["slug"])

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

	h.Repository.Save(vars["tenantID"], *p)
}

// Create godoc
// @Summary      Updates a Post
// @Description  Changes post body and updates LastUpdateTime timestamp
// @Accept       json
// @Produce      json
// @Param        tenantID   path      int  true  "Tenant ID"
// @Param        slug   path      string  true  "Unique slug of the post"
// @Param        {object} body UpdateRequest true "Post to create"
// @Success      200
// @Failure      400  {object}  responsehandler.Error
// @Failure      500  {object}  responsehandler.Error
// @Router       /v2/tenant/{tenantID}/posts/{slug} [post]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, err := h.Repository.GetBySlug(vars["tenantID"], vars["slug"])

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
		return
	}

	if p == nil {
		responsehandler.EncodeJSONError(w, nil, http.StatusNotFound)
		return
	}

	requestContents, err := io.ReadAll(r.Body)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
		return
	}

	var ur UpdateRequest
	json.Unmarshal(requestContents, &ur)

	p.UpdatedAt = time.Now()
	p.ContentRaw = ur.Body

	if err := h.Repository.Save(vars["tenantID"], *p); err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusInternalServerError)
		return
	}

	responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
}

// Create godoc
// @Summary      Delete a Post
// @Description  Completely removes a post
// @Accept       json
// @Produce      json
// @Param        tenantID   path      int  true  "Tenant ID"
// @Param        slug path string true "Unique slug of the post to delete"
// @Success      200
// @Failure      404  {object}  responsehandler.Error
// @Failure      502  {object}  responsehandler.Error
// @Router       /v2/tenant/{tenantID}/posts/{slug} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, err := h.Repository.GetBySlug(vars["tenantID"], vars["slug"])

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadGateway)
		return
	}

	if p == nil {
		responsehandler.EncodeJSONError(w, nil, http.StatusNotFound)
		return
	}

	h.Repository.DeleteByID(vars["tenantID"], p.ID)
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
		responsehandler.EncodeJSONError(w, fmt.Errorf("invalid FROM %d smaller than 0 or SIZE %d larger than 100", from_int, size_int), http.StatusBadRequest)
	}

	showDrafts, _ := strconv.ParseBool(vars["showDrafts"])

	filters := make(map[string]interface{})
	filters["IsPublished"] = true

	if showDrafts {
		filters["AuthorID"] = "abaltra"
		filters["IsPublished"] = nil
	}

	p, err := h.Repository.List(vars["tenantID"], from_int, size_int, filters)

	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusBadRequest)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}

// Create godoc
// @Summary      Retrieve a Post
// @Description  Retrieve an existing post
// @Accept       json
// @Produce      json
// @Param        tenantID   path      int  true  "Tenant ID"
// @Param        slug path string true "Unique slug of post to retrieve"
// @Success      200  {object}  post.Post
// @Failure      404  {object}  responsehandler.Error
// @Router       /v2/tenant/{tenantID}/posts/{slug} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	p, err := h.Repository.GetBySlug(vars["tenantID"], vars["slug"])
	if err != nil {
		responsehandler.EncodeJSONError(w, err, http.StatusNotFound)
	} else {
		responsehandler.EncodeJSONResponse(w, p, http.StatusOK, nil)
	}
}
