package main

import (
	"net/http"

	"lenslocked.com/controllers"

	"github.com/gorilla/mux"
)

func main() {
	usersC := controllers.NewUsers()
	staticC := controllers.NewStatic()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods(("GET"))
	r.Handle("/faq", staticC.Faq).Methods(("GET"))

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	http.ListenAndServe(":3000", r)
}
