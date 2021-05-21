package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to my super awesome site</h1>")
}

func contact(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:suppoer@lenslockedlcom\">support@lenslock.com</a>")

}

func faq(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<div><h1>FAQ</h1><p>What is this going on?</p></div>")
}

func notFound() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<h1>Not Found: 404</h1>")
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.NotFoundHandler = notFound()
	http.ListenAndServe(":3000", r)
}

