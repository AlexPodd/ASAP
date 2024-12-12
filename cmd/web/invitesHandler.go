package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AlexPodd/ASAP/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) myInvites(w http.ResponseWriter, r *http.Request) {
	Invites, err := app.invites.GetAllUserInvite(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Invites = Invites
	app.render(w, r, http.StatusOK, "myInvites.html", data)
}

func (app *application) acceptInvite(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	CompanyID, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || CompanyID < 1 {
		app.errorLog.Print(err)
		app.notFound(w)
		return
	}

	err = app.usersincompanies.Insert(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), CompanyID, "worker")
	if err != nil {
		if errors.Is(err, models.ErrInvalidUserID) {
			app.sessionManager.Put(r.Context(), "flash", "Invalid user!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			data := app.newTemplateData(r)
			app.render(w, r, http.StatusUnprocessableEntity, "home.html", data)
		} else if errors.Is(err, models.ErrDuplicateNameInCompany) {
			app.sessionManager.Put(r.Context(), "flash", "You are already in this company!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			data := app.newTemplateData(r)

			errDel := app.invites.DeleteInvite(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), CompanyID)
			if errDel != nil {
				app.serverError(w, err)
			}

			app.render(w, r, http.StatusUnprocessableEntity, "home.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	errDel := app.invites.DeleteInvite(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), CompanyID)
	if errDel != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have successfully joined the company!")
	app.sessionManager.Put(r.Context(), "flashtype", "info")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) declineInvite(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	CompanyID, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || CompanyID < 1 {
		app.errorLog.Print(err)
		app.notFound(w)
		return
	}

	errDel := app.invites.DeleteInvite(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), CompanyID)
	if errDel != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "You have successfully declined the invitation!")
	app.sessionManager.Put(r.Context(), "flashtype", "info")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) sendInviteMenu(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	UserToInviteID, err := strconv.Atoi(params.ByName("UserID"))
	if err != nil || UserToInviteID < 1 {
		app.errorLog.Print(err)
		app.notFound(w)
		return
	}

	CompanyWithUsers, err := app.usersincompanies.GetAllCompanyWhereUserAdminOrOwner(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.CompanyWithUsers = CompanyWithUsers
	data.UserToInviteID = UserToInviteID
	app.render(w, r, http.StatusOK, "sendInvite.html", data)
}

func (app *application) sendInvitePost(w http.ResponseWriter, r *http.Request) {
	UserToInviteID, CompanyID, err := ParseInviteRequest(r)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}
	message := r.FormValue("message")

	err = app.invites.AddInvite(UserToInviteID, CompanyID, message)

	if err != nil {
		app.serverError(w, err)
	}
	app.sessionManager.Put(r.Context(), "flash", "You have successfully send invite!")
	app.sessionManager.Put(r.Context(), "flashtype", "info")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ParseInviteRequest(r *http.Request) (int, int, error) {
	var UserToInviteID int
	var CompanyID int

	params := httprouter.ParamsFromContext(r.Context())
	UserToInviteID, err := strconv.Atoi(params.ByName("UserToInviteID"))
	if err != nil || UserToInviteID < 1 {
		return UserToInviteID, CompanyID, err
	}

	CompanyID, err = strconv.Atoi(params.ByName("companyID"))
	if err != nil || CompanyID < 1 {
		return UserToInviteID, CompanyID, err
	}

	return UserToInviteID, CompanyID, err
}
