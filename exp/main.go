package main

import (
	"fmt"

	"lenslocked.com/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	us.DestructiveReset()

	// Create a user
	user := models.User{
		Name:  "sunil",
		Email: "s@y.com",
	}

	if err := us.Create(&user); err != nil {
		panic(err)
	}

	fmt.Println("----- Find User BY ID -----")
	foundUser, err := us.ByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)

	fmt.Println("----- UPDATE NAME  -----")
	user.Name = "CHANGED"
	err = us.Update(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated USER NAME")
	foundUser, err = us.ByEmail("s@y.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)

	fmt.Println("----- DELETEING A USER BY ID -----")
	err = us.Delete(1)
	if err != nil {
		panic(err)
	}
	foundUser, err = us.ByEmail("s@y.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)

}
