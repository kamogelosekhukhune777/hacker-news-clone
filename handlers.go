package main

import "net/http"

func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	a.model.Users.Get()
}
