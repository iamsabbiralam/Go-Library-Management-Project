package handler

import "net/http"

type EmailForm struct {
	Email	string
	Errors	map[string]string
}

func (h *Handler) forgotPassword(rw http.ResponseWriter, r *http.Request) {
	form := EmailForm{}
	if err:= h.templates.ExecuteTemplate(rw, "reset-password.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}