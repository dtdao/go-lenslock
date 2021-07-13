package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

type Users struct {
	NewView *views.View
	us *models.UserService
}

type SignupForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
	Name string `schema:"name"`
}

func NewUsers(us *models.UserService) *Users{
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us: us,
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

	user := models.User{
		Name: form.Name,
		Email: form.Email,
	}

	if err := u.us.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, form)

	// fmt.Fprintln(w, r.PostForm["email"])
	// fmt.Fprintln(w, r.PostForm["password"])

	// fmt.Fprintln(w, "This is a fake messgae. Presend that we created a user")
}