package infrastructure

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Router(srv Server) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/balance/create", srv.createBalance).Methods(http.MethodPost)
	s.HandleFunc("/balance/topUp", srv.topUpBalance).Methods(http.MethodPost)
	s.HandleFunc("/balance/writeOff", srv.writeOffFromBalance).Methods(http.MethodPost)
	s.HandleFunc("/balance/remove", srv.removeBalance).Methods(http.MethodPost)
	return initJSONResponse(r)
}

func initJSONResponse(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
