package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Bookings struct {
	ID int `db:"id"`
	UserID int `db:"user_id"`
	BookID int `db:"book_id"`
	StartTime time.Time `db:"start_time"`
	EndTime time.Time `db:"end_time"`
	Start_time string
	End_time string
	BookName string
}

type FormBookings struct {
	Id int
	Booking Bookings
	Errors map[string]string
}

type MyBookings struct {
	Booking []Bookings
	Offset	int
	Limit	int
	Total	int
	TotalPage	int
	Paginate	[]BookingPagination
	CurrentPage	int
	NextPageURL	string
	PreviousPageURL	string
}

type BookingPagination struct {
	URL	string
	PageNumber	int
}

func (b *Bookings) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Start_time,
			validation.Required.Error("The Start Time Field is Required"),
		),
		validation.Field(&b.End_time,
			validation.Required.Error("The End Time Field is Required"),
		),
	)
}

func (h *Handler) createBookings(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	const getbook = "SELECT * FROM books WHERE id = $1"
	var book Book
	h.db.Get(&book, getbook, i)
	
	vErrs := map[string]string{}
	booking := Bookings{}
	h.loadCreateBookingForm(rw, i, booking, vErrs)
}

func(h *Handler) storeBookings(rw http.ResponseWriter, r *http.Request) {
	if err:= r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var booking Bookings
	if err := h.decoder.Decode(&booking, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := booking.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range vErrors {
				vErrs[key] = value.Error()
			}
			h.loadCreateBookingForm(rw, booking.ID, booking, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	const insertBooking = `INSERT INTO bookings(user_id,book_id,Start_time,end_time) VALUES($1,$2,$3,$4)`
	res:= h.db.MustExec(insertBooking, 1, booking.BookID, booking.Start_time, booking.End_time)
	getBook:= `UPDATE books SET status = false WHERE id = $1`
	h.db.MustExec(getBook, booking.BookID)

	if ok , err:= res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/book/list", http.StatusTemporaryRedirect)
}

func (h *Handler) loadCreateBookingForm(rw http.ResponseWriter, id int, booking Bookings, errs map[string]string) {
	form := FormBookings{
		Id: id,
		Booking: booking,
		Errors: errs,
	}
	if err:= h.templates.ExecuteTemplate(rw, "create-bookings.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func(h *Handler) myBookings(rw http.ResponseWriter, r *http.Request) {
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
	booking := []Bookings{}
	offset := 0
	limit := 4
	if p > 0 {
		offset = limit * p - limit
	}
	h.db.Select(&booking, "SELECT * FROM bookings offset $1 limit $2", offset, limit)
	total := 0
	h.db.Get(&total, "SELECT count(*) FROM bookings")
	nextPageURL := ""
	previousPageURL := ""
	totalPage := int(math.Ceil(float64(total)/float64(limit)))
	bookingPaginate := make([]BookingPagination, totalPage)
	for i := 0; i < totalPage; i++ {
		bookingPaginate[i] = BookingPagination{
			URL: fmt.Sprintf("http://localhost:3000/mybookings?page=%d", i + 1),
			PageNumber: i + 1,
		}
		if i + 1 == p {
			if i != 0 {
				previousPageURL = fmt.Sprintf("http://localhost:3000/mybookings?page=%d", i)
			}
			if i + 1 != totalPage {
				nextPageURL = fmt.Sprintf("http://localhost:3000/mybookings?page=%d", i + 2)
			}
		}
	}
	for key, value := range booking {
		const getBook = `SELECT book_name FROM books WHERE id = $1`
		var book Book
		h.db.Get(&book, getBook, value.BookID)
		start_time:= value.StartTime.Format("Mon Jan _2 2006 15:04 AM")
		end_time:= value.EndTime.Format("Mon Jan _2 2006 15:04 AM")
		booking[key].BookName = book.Book_name
		booking[key].Start_time = start_time
		booking[key].End_time = end_time
	}
	list := MyBookings{
		Booking: booking,
		Offset: offset,
		Limit: limit,
		Total: total,
		TotalPage: totalPage,
		Paginate: bookingPaginate,
		CurrentPage: p,
		NextPageURL: nextPageURL,
		PreviousPageURL: previousPageURL,
	}
	if err:= h.templates.ExecuteTemplate(rw, "my-bookings.html", list); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}