package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abaltra/blog/server/config"
	"github.com/abaltra/blog/server/post"
	"github.com/gorilla/mux"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	config := config.NewConfig()
	pm := &post.Repository{
		Config: config,
	}

	pm.Init()
	pm.Ping()

	ph := &post.Handler{
		Repository: pm,
	}

	router := mux.NewRouter()

	router.HandleFunc("/test", testHandler)
	router.HandleFunc("/posts", ph.List).Methods(http.MethodGet)
	router.HandleFunc("/posts", ph.Create).Methods(http.MethodPost)
	router.HandleFunc("/posts/{slug}", ph.Get).Methods(http.MethodGet)
	router.HandleFunc("/posts/{slug}", ph.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/posts/{slug}", ph.Update).Methods(http.MethodPost)
	router.HandleFunc("/posts/{slug}/publish", ph.Publish).Methods(http.MethodPut)

	srv := &http.Server{
		Handler:      router,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		Addr:         fmt.Sprintf("127.0.0.1:%s", config.Port),
	}

	fmt.Printf("Serving on 127.0.0.1:%s", config.Port)
	srv.ListenAndServe()
}
