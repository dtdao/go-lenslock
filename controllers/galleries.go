package  controllers

import (
	"lenslocked.com/models"
	"lenslocked.com/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New:   views.NewView("bootstrap", "galleries/new"),
		gs:        gs,
	}
}

type Galleries struct {
	New *views.View
	gs models.GalleryService
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
