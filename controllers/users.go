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

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrorNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrorInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := http.Cookie{
		Name: "email",
		Value: user.Email,
	}

	http.SetCookie(w, &cookie)
	fmt.Fprintln(w, user)
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Email is:", cookie.Value)
	fmt.Fprintln(w, cookie)
}