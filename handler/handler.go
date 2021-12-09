package handler

import (
	
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

const sessionName = "library-session"

type Handler struct {
	templates *template.Template
	db 	*sqlx.DB
	decoder *schema.Decoder
	sess *sessions.CookieStore
}

func New(db *sqlx.DB, decoder *schema.Decoder, sess *sessions.CookieStore) *mux.Router {
	h:= &Handler{
		db: db,
		decoder: decoder,
		sess: sess,
	}

	h.parseTemplate()

	r:= mux.NewRouter()
	r.HandleFunc("/", h.home)
	r.HandleFunc("/logout", h.logout)
	r.HandleFunc("/resetpassword", h.forgotPassword)

	l := r.NewRoute().Subrouter()
	l.HandleFunc("/registration", h.signUp).Methods("GET")
	l.HandleFunc("/registration", h.signUpCheck).Methods("POST")
	l.HandleFunc("/login", h.login).Methods("GET")
	l.HandleFunc("/login", h.loginCheck).Methods("POST")
	l.Use(h.loginMiddleware)

	s := r.NewRoute().Subrouter()
	s.Use(h.authMiddleware)
	s.HandleFunc("/category/create", h.createCategories)
	s.HandleFunc("/category/store", h.storeCategories)
	s.HandleFunc("/category/list", h.listCategories)
	s.HandleFunc("/category/{id:[0-9]+}/edit", h.editCategories)
	s.HandleFunc("/category/{id:[0-9]+}/update", h.updateCategories)
	s.HandleFunc("/category/{id:[0-9]+}/delete", h.deleteCategories)
	s.HandleFunc("/category/search", h.searchCategory)
	s.HandleFunc("/book/create", h.createBooks)
	s.HandleFunc("/book/store", h.storeBooks)
	s.HandleFunc("/book/list", h.listBooks)
	s.HandleFunc("/book/{id:[0-9]+}/edit", h.editBook)
	s.HandleFunc("/book/{id:[0-9]+}/update", h.updateBook)
	s.HandleFunc("/book/{id:[0-9]+}/delete", h.deleteBook)
	s.HandleFunc("/book/search", h.searchBook)
	s.HandleFunc("/bookings/{id:[0-9]+}/create", h.createBookings)
	s.HandleFunc("/bookings/store", h.storeBookings)
	s.HandleFunc("/mybookings", h.myBookings)
	s.HandleFunc("/book/{id:[0-9]+}/bookdetails", h.bookDetails)
	s.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./"))))
	

	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := h.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, "invalid URL", http.StatusInternalServerError)
			return
		}
	})

	return r
}

func (h *Handler) parseTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"templates/category/create-category.html",
		"templates/category/list-category.html",
		"templates/category/edit-category.html",
		"templates/category/404.html",
		"templates/book/create-book.html",
		"templates/book/list-book.html",
		"templates/book/edit-book.html",
		"templates/home.html",
		"templates/bookings/create-bookings.html",
		"templates/bookings/my-bookings.html",
		"templates/book/single-details.html",
		"templates/signup.html",
		"templates/login.html",
		"templates/reset-password.html",
		))
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, err := h.sess.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}
		authUserID := session.Values["authUserID"]
		if authUserID != nil {
			next.ServeHTTP(rw, r)
		} else {
			http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
		}
		
	})
}

func (h *Handler) loginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, err := h.sess.Get(r, sessionName)
		if err != nil {
			log.Fatal(err)
		}
		authUserID := session.Values["authUserID"]
		if authUserID != nil {
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}