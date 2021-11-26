package controllers

import (
	"fmt"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
	"log"
	"net/http"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

type GalleryForm struct {
	Title string `schema:"title"`
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	fmt.Println("create got the user:", user)
	gallery := models.Gallery{
		Title: form.Title,
		UserId: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	fmt.Fprintln(w, gallery)
}

//package controllers
//
//import (
//	"net/http"
//
//	"lenslocked.com/views"
//)
//
//type Galleries struct {
//	NewGallery *views.View
//}
//
//func NewGallery() *Galleries {
//	return &Galleries{
//		NewGallery: views.NewView("bootstrap", "galleries/new"),
//	}
//}
//
//func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
//	g.NewGallery.Render(w, nil)
//}
