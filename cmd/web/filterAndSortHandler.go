package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlexPodd/ASAP/internal/models"
)

// Пример фильтрации (замените на реальные данные или логику фильтрации)
func (app *application) filter(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	columnName := r.URL.Query().Get("column")
	paramName := r.URL.Query().Get("param")

	switch tableName {
	case "projects":
		companyID, err := strconv.Atoi(r.URL.Query().Get("companyID"))
		if err != nil || companyID < 1 {
			app.notFound(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		projects, err := filterProject(app, columnName, paramName, companyID)
		if err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(projects); err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
	case "ProjectsTask":
		companyID, err := strconv.Atoi(r.URL.Query().Get("companyID"))
		if err != nil || companyID < 1 {
			app.notFound(w)
			return
		}

		projectName := r.URL.Query().Get("projectName")

		projectID, err := app.projects.GetIDForName(projectName, companyID)
		if err != nil || projectID < 1 {
			app.notFound(w)
			return
		}
		app.infoLog.Print(tableName + " " + columnName + " " + paramName)
		w.Header().Set("Content-Type", "application/json")
		tasks, err := filterTask(app, columnName, paramName, companyID, projectID)
		if err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
	}
}
func (app *application) sort(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	columnName := r.URL.Query().Get("column")
	paramName := r.URL.Query().Get("param")

	switch tableName {
	case "company":
		w.Header().Set("Content-Type", "application/json")
		companies, err := sortCompany(app, columnName, paramName, r)
		if err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(companies); err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
	case "userInCompany":
		companyID, err := strconv.Atoi(r.URL.Query().Get("companyID"))
		if err != nil || companyID < 1 {
			app.notFound(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		users, err := sortUserInCompanies(app, columnName, paramName, companyID)
		if err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
	case "ProjectsTask":
		companyID, err := strconv.Atoi(r.URL.Query().Get("companyID"))
		if err != nil || companyID < 1 {
			app.notFound(w)
			return
		}
		projectName := r.URL.Query().Get("projectName")

		projectID, err := app.projects.GetIDForName(projectName, companyID)
		if err != nil || projectID < 1 {
			app.notFound(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		tasks, err := sortTask(app, columnName, paramName, companyID, projectID)
		if err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		}
	}

}

func filterProject(app *application, columnName, paramName string, companyID int) ([]*models.Project, error) {
	if columnName == "status" && paramName == "Completed project" {
		Projects, err := app.projects.GetAllCompanyProjectsFilteComplited(companyID)
		return Projects, err
	}
	if columnName == "status" && paramName == "Outstanding project" {
		Projects, err := app.projects.GetAllCompanyProjectsFilterOutstanding(companyID)
		return Projects, err
	}

	return nil, nil
}
func sortCompany(app *application, columnName, paramName string, r *http.Request) ([]*models.CompanyWithUsers, error) {
	if columnName == "role" && paramName == "ascending" {
		Companyies, err := app.usersincompanies.GetAllCompanyWithUserSortByRoleAscending(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
		return Companyies, err
	}
	if columnName == "role" && paramName == "descending" {
		Companyies, err := app.usersincompanies.GetAllCompanyWithUserSortByRoleDescending(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
		return Companyies, err
	}
	return nil, nil
}

func sortUserInCompanies(app *application, columnName, paramName string, companyID int) ([]*models.UserWithRole, error) {
	if columnName == "role" && paramName == "ascending" {
		Users, err := app.usersincompanies.GetAllUsersSortByRoleAscending(companyID)
		return Users, err
	}
	if columnName == "role" && paramName == "descending" {
		Users, err := app.usersincompanies.GetAllUsersSortByRoleDescending(companyID)
		return Users, err
	}
	return nil, nil
}

func sortTask(app *application, columnName, paramName string, companyID, projectID int) ([]*models.Task, error) {
	if columnName == "expired" && paramName == "ascending" {
		Tasks, err := app.tasks.GetAllCompanyProjectTasksSortByExpiredAscending(companyID, projectID)
		return Tasks, err
	}
	if columnName == "expired" && paramName == "descending" {
		Tasks, err := app.tasks.GetAllCompanyProjectTasksSortByExpiredDescending(companyID, projectID)
		return Tasks, err
	}

	if columnName == "category" && paramName == "ascending" {
		Tasks, err := app.tasks.GetAllCompanyProjectTasksSortByCategoryAscending(companyID, projectID)
		return Tasks, err
	}
	if columnName == "category" && paramName == "descending" {
		Tasks, err := app.tasks.GetAllCompanyProjectTasksSortByCategoryDescending(companyID, projectID)
		return Tasks, err
	}

	return nil, nil
}

func filterTask(app *application, columnName, paramName string, companyID, projectID int) ([]*models.Task, error) {
	app.infoLog.Print(columnName + " " + paramName)
	if columnName == "status" && paramName == "Completed task" {
		Tasks, err := app.tasks.GetAllCompanyProjectTasksComplited(companyID, projectID)
		app.errorLog.Print(Tasks)
		app.errorLog.Print(err)
		return Tasks, err
	}
	if columnName == "status" && paramName == "Outstanding task" {
		Tasks, err := app.tasks.GetAllCompanyProjectTasksUnomplited(companyID, projectID)
		app.errorLog.Print(Tasks)
		app.errorLog.Print(err)
		return Tasks, err
	}

	return nil, nil
}
