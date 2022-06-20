package main

import (
	"fmt"
	"glog/config"
	"glog/post"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	fmt.Println("Starting glog")
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
	router.HandleFunc("/tenant/{tenantID}/posts", ph.List).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID/posts", ph.Create).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/posts/{slug}", ph.Get).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/posts/{slug}", ph.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/tenant/{tenantID}/posts/{slug}", ph.Update).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/posts/{slug}/publish", ph.Publish).Methods(http.MethodPut)

	srv := &http.Server{
		Handler:      router,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		Addr:         fmt.Sprintf("127.0.0.1:%s", config.Port),
	}

	fmt.Printf("Serving on 127.0.0.1:%s", config.Port)
	srv.ListenAndServe()
}
