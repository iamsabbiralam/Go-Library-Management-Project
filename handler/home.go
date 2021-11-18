package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) home(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(rw, "<h1>Hello!</h1>")
}