package main

import (
	"fmt"
	"lenslocked.com/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dong"
	password = "password"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	// Todo fix this
	//defer us.Close()
	//us.DestructiveReset();
	//us.AutoMigrate()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	r := mux.NewRouter()
	galleriesC := controllers.NewGalleries(services.GalleryService, r)
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	//userByAge, err := services.InAgeRange(99, 100);
	if err != nil {
		panic(err)
	}
	//fmt.Println(userByAge)

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")

	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name("show_gallery")
	r.Handle("/login", usersC.LoginView).Methods("GET")

	// r.HandleFunc("/faq", faq).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+/edit}", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+/update}", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	http.ListenAndServe(":3000", r)
}

// func must(err error){
// 	if err != nil {
// 		panic(err)
// 	}
// }
