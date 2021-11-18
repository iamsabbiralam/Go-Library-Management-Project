package handler

import (
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
	r.HandleFunc("/category/create", h.createCategory)
	r.HandleFunc("/category/store", h.storeCategory)
	r.HandleFunc("/category/list", h.listCategories)

	return r
}

func (h *Handler) parseTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"templates/create-category.html",
		"templates/list-category.html",
		))
}