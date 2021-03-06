package models

import "gorm.io/gorm"

type Gallery struct {
	gorm.Model
	UserId uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}
type GalleryDB interface {
	Create(gallery *Gallery) error
	ById(id uint) (*Gallery, error)
	Update(gallery *Gallery) error
	Delete(id uint) error
}

type galleryService struct {
	GalleryDB
}
type galleryValidator struct {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db: db}},
	}
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) ById(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	return &gallery, err
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIdRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}
func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.userIdRequired,
		gv.titleRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}
func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.UserId <= 0 {
		return ErrTitleRequired
	}
	return nil
}
func (gv *galleryValidator) userIdRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrorIDInvalid
	}
	return gv.GalleryDB.Delete(id)
}
func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&gallery).Error
}

type galleryValFunc func(*Gallery) error

func runGalleryValFuncs(gallery *Gallery, fns ...galleryValFunc) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}
