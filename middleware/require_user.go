package middleware

import (
	"fmt"
	"net/http"

	"lenslocked.com/context"
	"lenslocked.com/models"
)

type RequireUser struct {
	models.UserService
}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user is logged in and then either
// call next.ServeHTTP(w, r) if they are, or redirect
// them to the login page if they are not.
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user is logged in and then either
// call next(w, r) if they are, or redirect them to
// the login page if they are not
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	// We want to return a dynamically created func(w, r) but we
	// also need to convert to http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a user is logged in
		// if yes, call next(w,r)
		// if not, http.redirect to "/login"
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		fmt.Println("User Found: ", user)

		// Get the context from our request
		ctx := r.Context()

		// create a new context from the existing one that has
		// our user stored in it with the private user key
		ctx = context.WithUser(ctx, user)

		// create a new request from the existing one with our context
		// attached to it and assign it back to 'r'
		r = r.WithContext(ctx)

		// call the next handler in the chain with updated context
		next(w, r)
	})
}
