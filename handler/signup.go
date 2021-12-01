package handler

import (
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type SignUp struct {
	ID	int `db:"id"`
	FirstName string `db:"first_name"`
	LastName string `db:"last_name"`
	Email string `db:"email"`
	Password string `db:"password"`
	ConfirmPassword string 
	IsVerified bool `db:"is_verified"`
}

type SignUpForm struct {
	SingUp	SignUp
	Errors	map[string]string
}

func (s *SignUp) Validate() error {
	return validation.ValidateStruct(s,
	validation.Field(&s.FirstName,
		validation.Required.Error("This field is must required")),
	validation.Field(&s.LastName,
		validation.Required.Error("This field is must required")),
	validation.Field(&s.Email,
		validation.Required.Error("This field is must required")),
	validation.Field(&s.Password,
		validation.Required.Error("This field is must required")),
	validation.Field(&s.ConfirmPassword,
		validation.Required.Error("This field is must required")))
} 

func (h *Handler) signUp(rw http.ResponseWriter, r *http.Request) {
	vErrs := map[string]string{}
	signup := SignUp{}
	h.loadSignUpForm(rw, signup, vErrs)
}

func (h *Handler) signUpCheck(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var signup SignUp
	if err := h.decoder.Decode(&signup, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if signup.Password != signup.ConfirmPassword {
		formData := SignUpForm{
			SingUp: signup,
			Errors: map[string]string{"Password" : "The password does not match with the confirm password"},
		}
		if err:= h.templates.ExecuteTemplate(rw, "signup.html", formData); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	}

	if err := signup.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range vErrors {
				vErrs[key] = value.Error()
				fmt.Println(key)
			}
			h.loadSignUpForm(rw, signup, vErrs)
			return
		}
	}

	const userSingUp = `INSERT INTO users(first_name, last_name, email, password) VALUES($1, $2, $3, $4)`
	pass, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	res := h.db.MustExec(userSingUp, signup.FirstName, signup.LastName, signup.Email, string(pass))
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}

func (h *Handler) loadSignUpForm(rw http.ResponseWriter, singup SignUp, errs map[string]string) {
	data := SignUpForm{
		SingUp: singup,
		Errors: errs,
	}
	if err:= h.templates.ExecuteTemplate(rw, "signup.html", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}