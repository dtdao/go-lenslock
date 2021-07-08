package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/views"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers() *Users{
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}


func (u *Users) New(w http.ResponseWriter, r *http.Request)  {
	u.NewView.Render(w, nil)
}


func (u *Users) Create(w http.ResponseWriter, r *http.Request){
	var form SignupForm
	if err := parseForm(r, &form);  err != nil {
		panic(err)
	}

	fmt.Fprintln(w, form)

	// fmt.Fprintln(w, r.PostForm["email"])
	// fmt.Fprintln(w, r.PostForm["password"])

	// fmt.Fprintln(w, "This is a fake messgae. Presend that we created a user")
}