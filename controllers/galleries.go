package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
	"log"
	"net/http"
	"strconv"
)

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		EditView: views.NewView("bootstrap", "galleries/edit"),
		gs:       gs,
		r:        r,
	}
}

type GalleryForm struct {
	Title string `schema:"title"`
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserId != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	var form GalleryForm
	vd.Yield = gallery

	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery successfully updated!",
	}
	g.EditView.Render(w, vd)
	fmt.Fprintln(w, gallery)
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserId != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

type Galleries struct {
	New      *views.View
	ShowView *views.View
	EditView *views.View
	gs       models.GalleryService
	r        *mux.Router
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
		Title:  form.Title,
		UserId: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	url, err := g.r.Get("show_gallery").URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
	fmt.Fprintln(w, gallery)
}
func (g *Galleries) galleryById(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Gallery Id", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ById(uint(id))
	if err != nil {
		switch err {
		case models.ErrorNotFound:
			http.Error(w, "Gallery not found ", http.StatusNotFound)
		default:
			http.Error(w, "Whoops. Something went wrong", http.StatusInternalServerError)
		}
		return nil, err
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
	return gallery, nil
}
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryById(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserId != user.ID {
		http.Error(w, "Gallery no found", http.StatusNotFound)
		return
	}
	var vd views.Data
	g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, vd)
		return
	}
	// TODO redirec to index
	fmt.Fprintln(w, "successfully deleted!")
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
