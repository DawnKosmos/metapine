package frontend

import "github.com/gorilla/mux"

func NewMux() *mux.Router {
	m := mux.NewRouter()
	return m
}
