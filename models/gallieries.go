package models

import "github.com/jinzhu/gorm"

// Gallery represents the galleries table in our DB
// and is mostly a container resource composed of images.
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

type GalleryService interface {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryValidator struct {
	GalleryDB
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
