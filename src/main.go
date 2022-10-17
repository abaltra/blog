package main

import (
	"fmt"
	"glog/config"
	"glog/post"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	_ "glog/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

// @title Glog - A Go Blogging backend using Mongo
// @version 1.0
// @description Very simple implementation of a bloggin platform using Mongo as a data store and Go with gorilla/mux.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host blog.abaltra.me/api
// @BasePath /v1
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

	router.PathPrefix("/swagger/*").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://127.0.0.1:%s/swagger/doc.json", config.Port)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	)).Methods(http.MethodGet)

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

	fmt.Printf("Serving on 127.0.0.1:%s\n", config.Port)
	log.Fatal(srv.ListenAndServe())
}
