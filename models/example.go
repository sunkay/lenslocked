package models

import "github.com/jinzhu/gorm"

// Gallery represents the galleries table in our DB
// and is mostly a container resource composed of images.
type Example struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

func NewExampleService(db *gorm.DB) ExampleService {
	return &exampleService{
		ExampleDB: &exampleValidator{
			ExampleDB: &exampleGorm{
				db: db,
			},
		},
	}
}

type ExampleService interface {
	ExampleDB
}

type exampleService struct {
	ExampleDB
}

type ExampleDB interface {
	Create(example *Example) error
}

type exampleValidator struct {
	ExampleDB
}

var _ ExampleDB = &exampleGorm{}

type exampleGorm struct {
	db *gorm.DB
}

func (gg *exampleGorm) Create(example *Example) error {
	return gg.db.Create(example).Error
}
