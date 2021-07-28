package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

type Users struct {
	NewView *views.View
	LoginView *views.View
	us *models.UserService
}

type SignupForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
	Name string `schema:"name"`
	Age uint `schema:"age"`
}

func NewUsers(us *models.UserService) *Users{
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
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
		Age: uint8(form.Age),
		Password: form.Password,
	}

	if err := u.us.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, user)

	// fmt.Fprintln(w, r.PostForm["email"])
	// fmt.Fprintln(w, r.PostForm["password"])

	// fmt.Fprintln(w, "This is a fake messgae. Presend that we created a user")
}

type LoginForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request){
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, form)
}