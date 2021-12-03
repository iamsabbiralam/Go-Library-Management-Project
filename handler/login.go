package handler

import (
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Email	string
	Password	string
	Errors	map[string]string
}

func (l *LoginForm) Validate() error {
	return validation.ValidateStruct(l,
	validation.Field(&l.Email,
		validation.Required.Error("The email field is must required")),
	validation.Field(&l.Password,
		validation.Required.Error("The password field is must required"),
		validation.Length(6, 32).Error("The password must be between 6 to 32 characters.")))
} 

func (h *Handler) login(rw http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err:= h.templates.ExecuteTemplate(rw, "login.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loginCheck(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var login LoginForm
	if err := h.decoder.Decode(&login, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := login.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range vErrors {
				vErrs[key] = value.Error()
			}
			login.Errors = vErrs
			if err:= h.templates.ExecuteTemplate(rw, "login.html", login); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
			}
		}
		return
	}

	userQuery := `SELECT * FROM users WHERE email = $1`
	var user SignUp
	h.db.Get(&user, userQuery, login.Email)
	if user.Email == "" {
		login.Errors = map[string]string{"Email" : "Invalid email given."}
		h.loadLoginForm(rw, login)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		login.Errors = map[string]string{"Password" : "Invalid password given."}
		h.loadLoginForm(rw, login)
		return
	}

	session, err := h.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	session.Options.HttpOnly = true

	session.Values["authUserID"] = user.ID
	if err := session.Save(r, rw); err != nil {
		log.Fatal(err)
	}
	if err := session.Save(r, rw); err != nil {
		log.Fatal(err)
	}

	http.Redirect(rw, r, "/book/list", http.StatusTemporaryRedirect)
}

func (h *Handler) loadLoginForm(rw http.ResponseWriter, login LoginForm) {
	if err:= h.templates.ExecuteTemplate(rw, "login.html", login); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}