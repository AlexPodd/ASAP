package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AlexPodd/ASAP/internal/models"
	"github.com/AlexPodd/ASAP/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type companyCreateForm struct {
	Name                string `form:"name"`
	validator.Validator `form:"-"`
}

func (app *application) getAllUsers(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	UsersWithRole, err := app.usersincompanies.GetAllUsers(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.CurrentCompanyId = id
	data.UsersWithRole = UsersWithRole
	app.render(w, r, http.StatusOK, "companyUsers.html", data)
}

func (app *application) getmyCompanies(w http.ResponseWriter, r *http.Request) {
	CompanyWithUsers, err := app.usersincompanies.GetAllCompanyWithUser(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.CompanyWithUsers = CompanyWithUsers
	app.render(w, r, http.StatusOK, "myCompanies.html", data)
}

func (app *application) getCompanyMenu(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	data := app.newTemplateData(r)
	data.CurrentCompanyId = id
	app.render(w, r, http.StatusOK, "companyMenu.html", data)
}

func (app *application) getCompanyProjects(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	CompaniesProjects, err := app.projects.GetAllCompanyProjects(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.CurrentCompanyId = id
	data.Projects = CompaniesProjects
	app.render(w, r, http.StatusOK, "companyProjects.html", data)
}
func (app *application) createCompany(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = companyCreateForm{}
	app.render(w, r, http.StatusOK, "createCompany.html", data)
}

func (app *application) createCompanyPost(w http.ResponseWriter, r *http.Request) {
	var form companyCreateForm
	var id int
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	id, err = app.company.Insert(form.Name, app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateCompanyName) {
			form.AddFieldError("name", "There is already a company with that name")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "createCompany.html", data)
		} else if errors.Is(err, models.ErrInvalidUserID) {
			form.AddFieldError("name", "Invalid user")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "createCompany.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.usersincompanies.Insert(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), id, "owner")
	app.sessionManager.Put(r.Context(), "flash", "Company creating was successful!")
	app.sessionManager.Put(r.Context(), "flashtype", "info")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) companyControl(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	UsersWithRole, err := app.usersincompanies.GetAllUsers(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.CurrentCompanyId = id
	data.UsersWithRole = UsersWithRole
	app.render(w, r, http.StatusOK, "companyControl.html", data)
}

func (app *application) companyControlPost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	companyID, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || companyID < 1 {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(params.ByName("userID"))
	if err != nil || userID < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	action := params.ByName("action")
	if action == "" {
		http.Error(w, "Missing action parameter", http.StatusBadRequest)
		return
	}

	switch action {
	case "makeAdmin":
		makeAdmin, err := app.usersincompanies.SetAdminRole(userID, companyID)
		if !makeAdmin || err != nil {
			app.sessionManager.Put(r.Context(), "flash", "An error occurred during make admin!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			app.sessionManager.Put(r.Context(), "flash", "Successful make admin!")
			app.sessionManager.Put(r.Context(), "flashtype", "info")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	case "makeUser":
		makeWorker, err := app.usersincompanies.SetWorkerRole(userID, companyID)
		if !makeWorker || err != nil {
			app.sessionManager.Put(r.Context(), "flash", "An error occurred make worker!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			app.sessionManager.Put(r.Context(), "flash", "Successful make worker!")
			app.sessionManager.Put(r.Context(), "flashtype", "info")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	case "remove":
		rem, err := app.usersincompanies.DeleteUser(userID, companyID)
		if !rem || err != nil {
			app.sessionManager.Put(r.Context(), "flash", "An error occurred during deletion!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			app.sessionManager.Put(r.Context(), "flash", "Successful removal!")
			app.sessionManager.Put(r.Context(), "flashtype", "info")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

	default:
		app.sessionManager.Put(r.Context(), "flash", "Invalid action!")
		app.sessionManager.Put(r.Context(), "flashtype", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

}
