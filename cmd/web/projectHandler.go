package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlexPodd/ASAP/internal/models"
	"github.com/AlexPodd/ASAP/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type projectCreateForm struct {
	Name                string `form:"name"`
	validator.Validator `form:"-"`
}

func (app *application) createProj(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	data := app.newTemplateData(r)
	data.Form = projectCreateForm{}
	data.CurrentCompanyId = id
	app.render(w, r, http.StatusOK, "createProject.html", data)

}

func (app *application) createProjPost(w http.ResponseWriter, r *http.Request) {
	var form projectCreateForm
	var id int

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.CurrentCompanyId = id
		app.render(w, r, http.StatusUnprocessableEntity, fmt.Sprintf("/company/projects/%d", id), data)
		return
	}

	err = app.projects.Insert(form.Name, id, app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidProjectName) {
			form.AddFieldError("name", "There is already a project with that name")
			data := app.newTemplateData(r)
			data.CurrentCompanyId = id
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "createProject.html", data)
		} else if errors.Is(err, models.ErrInvalidUserID) {
			form.AddFieldError("name", "Invalid user")
			data := app.newTemplateData(r)
			data.CurrentCompanyId = id
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "createProject.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.CurrentCompanyId = id

	app.sessionManager.Put(r.Context(), "flash", "Project creating was successful!")
	app.sessionManager.Put(r.Context(), "flashtype", "info")

	http.Redirect(w, r, fmt.Sprintf("/company/projects/%d", id), http.StatusSeeOther)
}
