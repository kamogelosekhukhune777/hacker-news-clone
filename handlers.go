package main

import "net/http"

//home page
func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	a.model.Users.Get()
}
