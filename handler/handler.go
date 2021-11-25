package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	templates *template.Template
	db 	*sqlx.DB
	decoder *schema.Decoder
}

func New(db *sqlx.DB, decoder *schema.Decoder) *mux.Router {
	h:= &Handler{
		db: db,
		decoder: decoder,
	}
	h.parseTemplate()

	r:= mux.NewRouter()
	r.HandleFunc("/", h.home)
	r.HandleFunc("/category/create", h.createCategories)
	r.HandleFunc("/category/store", h.storeCategories)
	r.HandleFunc("/category/list", h.listCategories)
	r.HandleFunc("/category/{id:[0-9]+}/edit", h.editCategories)
	r.HandleFunc("/category/{id:[0-9]+}/update", h.updateCategories)
	r.HandleFunc("/category/{id:[0-9]+}/delete", h.deleteCategories)
	r.HandleFunc("/category/search", h.searchCategory)
	r.HandleFunc("/book/create", h.createBooks)
	r.HandleFunc("/book/store", h.storeBooks)
	r.HandleFunc("/book/list", h.listBooks)
	r.HandleFunc("/book/{id:[0-9]+}/edit", h.editBook)
	r.HandleFunc("/book/{id:[0-9]+}/update", h.updateBook)
	r.HandleFunc("/book/{id:[0-9]+}/delete", h.deleteBook)
	r.HandleFunc("/book/search", h.searchBook)
	r.HandleFunc("/bookings/{id:[0-9]+}/create", h.createBookings)
	r.HandleFunc("/bookings/store", h.storeBookings)
	r.HandleFunc("/mybookings", h.myBookings)
	r.HandleFunc("/book/{id:[0-9]+}/bookdetails", h.bookDetails)

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
		))
}