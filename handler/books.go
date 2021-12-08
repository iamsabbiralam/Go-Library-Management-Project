package handler

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Book struct {
	ID int `db:"id"`
	Category_id int `db:"category_id"`
	Book_name string `db:"book_name"`
	AuthorName string `db:"author_name"`
	Details string `db:"details"`
	Image string `db:"image"`
	Status bool `db:"status"`
	Cat_name string
}

type FormBooks struct {
	Book Book
	Category []Category
	Errors map[string]string
}

type showBooks struct {
	Book	[]Book
	Booking	[]Bookings
	Category	[]Category
	Offset	int
	Limit	int
	Total	int
	Paginate	[]Pagination
	CurrentPage	int
	NextPageURL	string
	PreviousPageURL	string
	Search	string
}

type Pagination struct {
	URL	string
	PageNumber	int
}

func (b *Book) Validate() error {
	return validation.ValidateStruct(b, 
		validation.Field(&b.Book_name, 
			validation.Required.Error("This field is must be required"),
			validation.Length(3,0).Error("This field is must be grater than 3"),
		),
		validation.Field(&b.AuthorName,
			validation.Required.Error("The Author Name Field is Required"),
		),
		validation.Field(&b.Details,
			validation.Required.Error("The Details Field is Required"),
		))
}

func (h *Handler) createBooks(rw http.ResponseWriter, r *http.Request) {
	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")
	vErrs := map[string]string{}
	book := Book{}
	h.loadCreateBookForm(rw, book, category, vErrs)
}

func (h *Handler) storeBooks(rw http.ResponseWriter, r *http.Request) {
	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var book Book
	if err:= h.decoder.Decode(&book, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("Image")

	if file == nil {
		vErrs := map[string]string{"Image" : "The image field is required"}
		h.loadCreateBookForm(rw, book, category, vErrs)
			return
	}

    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
    defer file.Close()

    tempFile, err := ioutil.TempFile("assets/image", "upload-*.png")
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
    defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
	tempFile.Write(fileBytes)
	
	imageName := tempFile.Name()

	if err := book.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range vErrors {
				vErrs[key] = value.Error()
			}
			h.loadCreateBookForm(rw, book, category, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	const insertBook = `INSERT INTO books(category_id,book_name, author_name, details, image, status) VALUES($1, $2, $3, $4, $5, $6)`
	res:= h.db.MustExec(insertBook, book.Category_id, book.Book_name, book.AuthorName, book.Details, imageName, book.Status)
	if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/book/list", http.StatusTemporaryRedirect)
}

func(h *Handler) listBooks(rw http.ResponseWriter, r *http.Request) {
	// currentTime := time.Now()
	// booking := []Bookings{}
	// const getBooking = "SELECT * FROM bookings WHERE end_time < $1"
	// h.db.Select(&booking, getBooking, currentTime)
	// for _, value := range booking {
	// 	const updateBook = "UPDATE books SET status = true WHERE id = $1"
	// 	h.db.MustExec(updateBook, value.BookID)
	// }
	page := r.URL.Query().Get("page")
	var p int = 1
	var err error
	if page != "" {
		p, err = strconv.Atoi(page)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	book := []Book{}
	offset := 0
	limit := 3
	nextPageURL := ""
	previousPageURL := ""
	if p > 0 {
		offset = limit * p - limit
	}
	total := 0
	h.db.Get(&total, `SELECT count(*) FROM books`)
	h.db.Select(&book, "SELECT * FROM books offset $1 limit $2", offset, limit)
	for key, value := range book {
		const getTodo = `SELECT name FROM categories WHERE id=$1`
		var category Category
		h.db.Get(&category, getTodo, value.Category_id)
		book[key].Cat_name = category.Name
	}

	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")

	totalPage := int(math.Ceil(float64(total)/float64(limit)))

	paginate := make([]Pagination, totalPage)
	for i := 0; i < totalPage; i++ {
		paginate[i] = Pagination{
			URL: fmt.Sprintf("http://localhost:3000/book/list?page=%d", i + 1),
			PageNumber: i + 1,
		}
		if i + 1 == p {
			if i != 0 {
				previousPageURL = fmt.Sprintf("http://localhost:3000/book/list?page=%d", i)
			}
			if i + 1 != totalPage {
				nextPageURL = fmt.Sprintf("http://localhost:3000/book/list?page=%d", i + 2)
			}
		}
	}
	list := showBooks{
		Book : book,
		Category: category,
		Offset: offset,
		Limit: limit,
		Total: total,
		Paginate: paginate,
		CurrentPage: p,
		NextPageURL: nextPageURL,
		PreviousPageURL: previousPageURL,
	}

	if err:= h.templates.ExecuteTemplate(rw, "list-book.html", list); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) editBook(rw http.ResponseWriter, r *http.Request) {
	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}
	const getBook = `SELECT * FROM books WHERE id=$1`
	var book Book
	h.db.Get(&book, getBook, id)
	if book.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}
	h.loadEditBookForm(rw, book, category, map[string]string{})
}

func (h *Handler) updateBook(rw http.ResponseWriter, r *http.Request) {
	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}
	const getBook = `SELECT * FROM books WHERE id=$1`
	var book Book
	h.db.Get(&book, getBook, id)

	if book.ID == 0 {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	if err := h.decoder.Decode(&book, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("Image")
    
	var imageName string
	
    if err == nil {
		defer file.Close()
		tempFile, err := ioutil.TempFile("assets/image", "upload-*.png")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		
		tempFile.Write(fileBytes)
		
		imageName = tempFile.Name()

		if err := os.Remove(book.Image); err != nil {
				http.Error(rw, "Unable to upload image", http.StatusInternalServerError)
				return
			}
	} else {
		imageName = book.Image
	}

	if err := book.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range vErrors {
				vErrs[key] = value.Error()
			}
			h.loadEditBookForm(rw, book, category, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	const updateBook = `UPDATE books SET category_id = $2, book_name = $3, author_name = $4, details = $5, image = $6, status = $7 WHERE id = $1`
	res:= h.db.MustExec(updateBook, id, book.Category_id, book.Book_name, book.AuthorName, book.Details, imageName, book.Status)
	if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/book/list", http.StatusTemporaryRedirect)
}

func (h *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getbook = "SELECT * FROM books WHERE id = $1"
	var book Book
	h.db.Get(&book, getbook, id)

	if book.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const deleteBook = `DELETE FROM books WHERE id = $1`
	res:= h.db.MustExec(deleteBook, id)
	if ok, err:= res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/book/list", http.StatusTemporaryRedirect)
}

func (h *Handler) loadCreateBookForm(rw http.ResponseWriter, book Book, cat []Category, errs map[string]string) {
	form := FormBooks{
		Book : book,
		Category: cat,
		Errors : errs,
	}
	if err:= h.templates.ExecuteTemplate(rw, "create-book.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadEditBookForm(rw http.ResponseWriter, book Book, cat []Category, errs map[string]string) {
	form := FormBooks{
		Category : cat,
		Book : book,
		Errors : errs,
	}
	if err:= h.templates.ExecuteTemplate(rw, "edit-book.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) searchBook(rw http.ResponseWriter, r *http.Request) {
	if err:= r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	search := r.FormValue("search")
	const getSearch = "SELECT * FROM books WHERE book_name ILIKE '%%' || $1 || '%%'"
	book := []Book{}
	h.db.Select(&book, getSearch, search)
	for key, value := range book {
		const getTodo = `SELECT name FROM categories WHERE id=$1`
		var category Category
		h.db.Get(&category, getTodo, value.Category_id)
		book[key].Cat_name = category.Name
	}
	list := showBooks{
		Book : book,
		Search: search,
	}
	if err:= h.templates.ExecuteTemplate(rw, "list-book.html", list); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) bookDetails(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}
	const getBook = `SELECT * FROM books WHERE id=$1`
	var book Book
	h.db.Get(&book, getBook, id)
	const getTodo = `SELECT name FROM categories WHERE id=$1`
	var category Category
	h.db.Get(&category, getTodo, book.Category_id)
	book.Cat_name = category.Name

	if err:= h.templates.ExecuteTemplate(rw, "single-details.html", book); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}