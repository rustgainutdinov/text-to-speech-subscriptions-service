package infrastructure

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Router(srv Server) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/balance/create", srv.createBalance).Methods(http.MethodGet)
	return r
}
