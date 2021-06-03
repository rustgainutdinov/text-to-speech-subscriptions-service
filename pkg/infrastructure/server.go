package infrastructure

import (
	"fmt"
	"net/http"
)

type Server interface {
	createBalance(w http.ResponseWriter, r *http.Request)
}

type server struct {
	dependencyContainer DependencyContainer
}

func (s *server) createBalance(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Ok")
}

func NewServer(dependencyContainer DependencyContainer) Server {
	return &server{dependencyContainer: dependencyContainer}
}
