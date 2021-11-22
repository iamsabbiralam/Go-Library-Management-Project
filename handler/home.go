package handler

import (
	"net/http"
)

func (h *Handler) home(rw http.ResponseWriter, r *http.Request) {
		if err:= h.templates.ExecuteTemplate(rw, "home.html", nil); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}