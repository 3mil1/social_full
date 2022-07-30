package muxHandler

import "net/http"

type Handler struct {
	Mux *http.ServeMux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Mux.ServeHTTP(w, r)
}
