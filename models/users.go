package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"lenslocked.com/hash"
	"lenslocked.com/rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const hmacSecretKey = "my_secret-hmac-key"

var (
	// ErrNotFound is returned when a resource cannot be found in DB
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when passwordHash
	// does not match the incoming password
	ErrInvalidPassword = errors.New("models: incorrect password provided")

	userPwPepper = "lived in west ford"
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
}

// Creates a user in the database with the associated hashed password
// This function makes sure we do not store raw password and only stores hashed
// passwords with the user record
func (us *UserService) Create(user *User) error {

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password+userPwPepper),
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)

	return us.db.Create(user).Error
}

// ByID will look up a user with the provided ID.
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// As a general rule, any error but ErrNotFound should
// probably result in a 500 error.
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail will look up a user with the provided Email Address.
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// As a general rule, any error but ErrNotFound should
// probably result in a 500 error.
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByRemember will lookup & return a user record that matches the provided Remember token.
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	db := us.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update will update the provided user with all of the data
// in the provided user object.
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}

	return us.db.Save(user).Error
}

// Delete will delete the user with provided ID
// Delete will return an ErrInvalidID if ID provided is 0
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(user).Error
}

// Authenticate a user by comparing the input email & password with
// the stored users hashed password. Returns User and Error
// If it is a match return foundUser, nil
// If email not found rerturn nil, ErrNotFound
// If password does not match return nil, ErrInvalidPassword
// otherwise return nil, error
func (us *UserService) Authenticate(email string, pwd string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(pwd+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}

}

// AutoMigrate will attempt to automatically migrate the Users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// DestructiveReset will drop all our tables and resets the database
// This should not be used normally, but will help when writing tests
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

func (us *UserService) Close() error {
	return us.db.Close()
}

// first will query the database supplied by gorm.DB and place the first
// record returned in dst. If nothing is found it will return ErrNotFound

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
