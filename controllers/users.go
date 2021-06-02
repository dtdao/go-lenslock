package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
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
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request)  {
	u.NewView.Render(w, nil)
}


func (u *Users) Create(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	dec := schema.NewDecoder()

	var form SignupForm

	if err := dec.Decode(&form, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)

	// fmt.Fprintln(w, r.PostForm["email"])
	// fmt.Fprintln(w, r.PostForm["password"])

	fmt.Fprintln(w, "This is a fake messgae. Presend that we created a user")
}