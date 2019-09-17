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

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()

	// Create a user
	user := models.User{
		Name:     "test_user",
		Email:    "test@test.com",
		Password: "test123",
	}

	if err := services.User.Create(&user); err != nil {
		panic(err)
	}

	//verify if remember token is present
	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Users remember token is missing...")
	}

	fmt.Println("----- Find User BY Remember Token -----")
	foundUser, err := services.User.ByRemember(user.Remember)
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)

	/*
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
	*/
}
