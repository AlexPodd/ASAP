package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) usersTable(w http.ResponseWriter, r *http.Request) {
	Users, err := app.users.UsersTable()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Users = Users
	app.render(w, r, http.StatusOK, "usersTable.html", data)
}

func (app *application) adminDeleteUser(w http.ResponseWriter, r *http.Request) {
	adminID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	params := httprouter.ParamsFromContext(r.Context())

	UserID, err := strconv.Atoi(params.ByName("userID"))
	if err != nil || UserID < 1 {
		app.notFound(w)
		return
	}
	if adminID == UserID {
		app.sessionManager.Put(r.Context(), "flash", "You cannot delete yourself!")
		app.sessionManager.Put(r.Context(), "flashtype", "error")
		http.Redirect(w, r, "/admin/usersTable", http.StatusSeeOther)
		return
	}

	rem, err := app.users.DeleteUser(UserID)

	if !rem || err != nil {
		app.sessionManager.Put(r.Context(), "flash", "An error occurred while deleting a user!")
		app.sessionManager.Put(r.Context(), "flashtype", "error")
		app.errorLog.Print(err)
		http.Redirect(w, r, "/admin/usersTable", http.StatusSeeOther)
	} else {
		app.sessionManager.Put(r.Context(), "flash", "You have successfully deleted the user!")
		app.sessionManager.Put(r.Context(), "flashtype", "info")
		http.Redirect(w, r, "/admin/usersTable", http.StatusSeeOther)
	}
}

func (app *application) companyTable(w http.ResponseWriter, r *http.Request) {
	Comapny, err := app.company.CompanyTable()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Company = Comapny
	app.render(w, r, http.StatusOK, "companyTable.html", data)
}

func (app *application) adminDeleteCompany(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	companyID, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || companyID < 1 {
		app.notFound(w)
		return
	}

	rem, err := app.company.DeleteCompany(companyID)

	if !rem || err != nil {
		app.sessionManager.Put(r.Context(), "flash", "An error occurred while deleting a company!")
		app.sessionManager.Put(r.Context(), "flashtype", "error")
		app.errorLog.Print(err)
		http.Redirect(w, r, "/admin/companyTable", http.StatusSeeOther)
	} else {
		app.sessionManager.Put(r.Context(), "flash", "You have successfully deleted the company!")
		app.sessionManager.Put(r.Context(), "flashtype", "info")
		http.Redirect(w, r, "/admin/companyTable", http.StatusSeeOther)
	}
}

func (app *application) adminPanel(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "adminpanel.html", data)
}
