package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/api"
	"github.com/temelpa/timetravel/service"
	"github.com/temelpa/timetravel/storage"
)

// logError logs all non-nil errors
func logError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func main() {
	router := mux.NewRouter()
	db, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	memService := service.NewInMemoryRecordService()
	dbService := service.NewDatabaseService(db)
	api := api.NewAPI(&memService, dbService)

	apiRoute := router.PathPrefix("/api/v1").Subrouter()
	apiRouteV2 := router.PathPrefix("/api/v2").Subrouter()
	apiRoute.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})
	api.CreateRoutes(apiRoute)
	api.CreateRoutesV2(apiRouteV2)

	address := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening on %s", address)
	log.Fatal(srv.ListenAndServe())
}
