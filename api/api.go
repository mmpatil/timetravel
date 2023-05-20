package api

import (
	"github.com/gorilla/mux"
	"github.com/temelpa/timetravel/service"
)

type API struct {
	records   service.RecordService
	recordsV2 service.DatabaseService
}

func NewAPI(records service.RecordService, recordsV2 service.DatabaseService) *API {
	return &API{records, recordsV2}
}

// generates all api routes
func (a *API) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}

// generates all api routes
func (a *API) CreateRoutesV2(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecordsV2).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecordsV2).Methods("POST")
}
