package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Yo! Welcome to Go World</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, `To get in touch please send an email to 
		<a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>`)
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
