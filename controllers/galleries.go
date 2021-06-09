package controllers

import (
	"net/http"

	"lenslocked.com/views"
)

type Galleries struct {
	NewGallery *views.View
}

func NewGallery() *Galleries {
	return &Galleries{
		NewGallery: views.NewView("bootstrap", "galleries/new"),
	}
}

func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewGallery.Render(w, nil)
}