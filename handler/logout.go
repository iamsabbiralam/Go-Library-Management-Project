
  
package handler

import (
	"log"
	"net/http"
)

func (h *Handler) logout(rw http.ResponseWriter, r *http.Request) {
	session, err := h.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}
	session.Values["authUserID"] = nil
	if err := session.Save(r, rw); err != nil {
		log.Fatal(err)
	}

	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}