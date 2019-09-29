package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil,
			fmt.Errorf("Unable to open postgres conn.. actual error: %s", err.Error())
	}
	db.LogMode(true)

	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		db:      db,
	}, nil
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

// DestructiveReset will drop all our tables and resets the database
// This should not be used normally, but will help when writing tests
func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}
