package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AlexPodd/ASAP/internal/models"
	"github.com/AlexPodd/ASAP/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type taskCreateForm struct {
	Name                string `form:"name"`
	Category            string `form:"category"`
	Expires             string `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	idComp, projectName, err := ParseTaskRequest(r)
	if err != nil {
		app.notFound(w)
		return
	}
	data := GenerateTaskTemplateData(app, r, idComp, projectName)
	data.Form = taskCreateForm{}
	app.render(w, r, http.StatusOK, "createTask.html", &data)
}

func (app *application) createTaskPost(w http.ResponseWriter, r *http.Request) {
	var form taskCreateForm

	idComp, projectName, err := ParseTaskRequest(r)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	TimeExpires, err := ParseExpires(form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	data := GenerateTaskTemplateData(app, r, idComp, projectName)

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Category), "category", "This field cannot be blank")
	//form.CheckField(TimeExpires.IsZero(), "expires", "Expiration date is required")

	if !form.Valid() {
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "createTask.html", &data)
		return
	}

	projID, err := app.projects.GetIDForName(projectName, idComp)
	if err != nil {
		app.clientError(w, 400)
	}

	err = app.tasks.Insert(form.Name, form.Category, TimeExpires, projID, idComp)
	if err != nil {
		if errors.Is(err, models.ErrInvalidTaskName) {
			form.AddFieldError("name", "There is already a task with that name")
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "createTask.html", &data)
		} else if errors.Is(err, models.ErrInvalidUserID) {
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Task created successfully!")
	app.sessionManager.Put(r.Context(), "flashtype", "info")
	http.Redirect(w, r, fmt.Sprintf("/tasks/%d/%s", idComp, projectName), http.StatusSeeOther)
}

func (app *application) viewAllTask(w http.ResponseWriter, r *http.Request) {
	idComp, projectName, err := ParseTaskRequest(r)
	if err != nil {
		app.notFound(w)
		return
	}
	data := GenerateTaskTemplateData(app, r, idComp, projectName)

	idProj, err := app.projects.GetIDForName(projectName, idComp)
	if err != nil {
		app.notFound(w)
		return
	}
	ProjectsTask, err := app.tasks.GetAllCompanyProjectTasks(idComp, idProj)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.ProjectsTask = ProjectsTask
	app.render(w, r, http.StatusOK, "projectTasks.html", &data)
}

func (app *application) completeTask(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idComp, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || idComp < 1 {
		app.notFound(w)
		return
	}
	projectName := params.ByName("projectName")
	taskName := params.ByName("taskName")
	idProj, err := app.projects.GetIDForName(projectName, idComp)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.tasks.CompleateTask(idProj, idComp, app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), taskName)
	if err != nil {
		if errors.Is(err, models.TaskNotFound) {
			app.sessionManager.Put(r.Context(), "flash", "Task not found")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, fmt.Sprintf("/tasks/%d/%s", idComp, projectName), http.StatusSeeOther)
		} else if errors.Is(err, models.TaskIsAlredyDone) {
			app.sessionManager.Put(r.Context(), "flash", "Task is alredy done")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, fmt.Sprintf("/tasks/%d/%s", idComp, projectName), http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You have successfully compleate the task!")
	app.sessionManager.Put(r.Context(), "flashtype", "info")
	http.Redirect(w, r, fmt.Sprintf("/tasks/%d/%s", idComp, projectName), http.StatusSeeOther)

}

func ParseExpires(form taskCreateForm) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", form.Expires)
}

func GenerateTaskTemplateData(app *application, r *http.Request, idComp int, projectName string) templateData {
	data := app.newTemplateData(r)
	data.CurrentCompanyId = idComp
	data.CurrentProjectName = projectName
	return *data
}

func ParseTaskRequest(r *http.Request) (int, string, error) {
	var idComp int
	var projectName string

	params := httprouter.ParamsFromContext(r.Context())
	idComp, err := strconv.Atoi(params.ByName("companyID"))
	if err != nil || idComp < 1 {
		return idComp, projectName, err
	}
	projectName = params.ByName("projectName")
	return idComp, projectName, err
}
