package handler

import (
	"log"
	"net/http"
)

type Auth struct {
	Auth	interface{}
}

func (h *Handler) home(rw http.ResponseWriter, r *http.Request) {
	session, err := h.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}
	auth := session.Values["authUserID"]
	list := Auth{
		Auth: auth,
	}
	if err:= h.templates.ExecuteTemplate(rw, "home.html", list); err != nil {
	http.Error(rw, err.Error(), http.StatusInternalServerError)
	return
	}
}