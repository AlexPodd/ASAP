package main

import (
	"net/http"
)

func (app *application) findUsers(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")

	if pattern == "" {
		app.sessionManager.Put(r.Context(), "flash", "Enter the data!")
		app.sessionManager.Put(r.Context(), "flashtype", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	Users, err := app.users.FindForIdOrUsername(pattern)

	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Users = Users
	data.Pattern = pattern
	app.render(w, r, http.StatusOK, "findUsers.html", data)

}
