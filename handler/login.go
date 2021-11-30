package handler

import "net/http"

func (h *Handler) login(rw http.ResponseWriter, r *http.Request) {
	if err:= h.templates.ExecuteTemplate(rw, "login.html", nil); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}