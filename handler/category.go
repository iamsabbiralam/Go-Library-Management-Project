package handler

import (
	"net/http"
)

type Category struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Status bool `db:"status"`
}

type FormCategory struct {
	Cat Category
	Errors map[string]string
}

type ListCategory struct {
	Categories []Category
}

func (h *Handler) createCategory(rw http.ResponseWriter, r *http.Request) {
	vErrs := map[string]string{}
	cat := Category{}
	h.loadCreateCategoryForm(rw, cat, vErrs)
}

func (h *Handler) storeCategory(rw http.ResponseWriter, r *http.Request) {
	if err:= r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if category.Name == "" {
		vErrs := map[string]string{
			"Name" : "This field is required",
		}
		h.loadCreateCategoryForm(rw, category, vErrs)
	}
	if len(category.Name) < 3 {
		vErrs := map[string]string{
			"Name" : "This field must be grater than 3",
		}
		h.loadCreateCategoryForm(rw, category, vErrs)
	}
	
	const insertCategory = `INSERT INTO category(name,status) VALUES($1,$2)`
	res:= h.db.MustExec(insertCategory, category.Name, category.Status )

	if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) listCategories(rw http.ResponseWriter, r *http.Request) {
	category := []Category{}
	h.db.Select(&category, "SELECT * FROM category")
	list := ListCategory{
		Categories: category,
	}
	if err:= h.templates.ExecuteTemplate(rw, "list-category.html", list); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadCreateCategoryForm(rw http.ResponseWriter, cat Category, errs map[string]string) {
	q := FormCategory{
		Cat : cat,
		Errors : errs,
	}
	if err:= h.templates.ExecuteTemplate(rw, "create-category.html", q); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}