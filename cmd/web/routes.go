package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Add the authenticate() middleware to the chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/company/create", protected.ThenFunc(app.createCompany))
	router.Handler(http.MethodPost, "/company/create", protected.ThenFunc(app.createCompanyPost))
	router.Handler(http.MethodGet, "/company/myCompanies/", protected.ThenFunc(app.getmyCompanies))
	router.Handler(http.MethodGet, "/findUsers", protected.ThenFunc(app.findUsers))
	router.Handler(http.MethodGet, "/myInvites", protected.ThenFunc(app.myInvites))
	router.Handler(http.MethodPost, "/acceptInvite/:companyID", protected.ThenFunc(app.acceptInvite))
	router.Handler(http.MethodPost, "/declineInvite/:companyID", protected.ThenFunc(app.declineInvite))
	router.Handler(http.MethodGet, "/sendInviteMenu/:UserID", protected.ThenFunc(app.sendInviteMenu))
	router.Handler(http.MethodGet, "/filter", protected.ThenFunc(app.filter))
	router.Handler(http.MethodGet, "/sort", protected.ThenFunc(app.sort))

	companyeis := protected.Append(app.isMember)

	router.Handler(http.MethodGet, "/company/users/:companyID", companyeis.ThenFunc(app.getAllUsers))
	router.Handler(http.MethodGet, "/company/menu/:companyID", companyeis.ThenFunc(app.getCompanyMenu))
	router.Handler(http.MethodGet, "/company/projects/:companyID", companyeis.ThenFunc(app.getCompanyProjects))
	router.Handler(http.MethodGet, "/tasks/:companyID/:projectName", companyeis.ThenFunc(app.viewAllTask))
	router.Handler(http.MethodGet, "/filter/:companyID", companyeis.ThenFunc(app.filter))
	router.Handler(http.MethodGet, "/sort/:companyID", companyeis.ThenFunc(app.sort))
	router.Handler(http.MethodGet, "/CompleteTask/:companyID/:projectName/:taskName", companyeis.ThenFunc(app.completeTask))

	adminInCompanyeis := protected.Append(app.isAdmin)

	router.Handler(http.MethodGet, "/task/create", adminInCompanyeis.ThenFunc(app.createTask))
	router.Handler(http.MethodPost, "/task/create", adminInCompanyeis.ThenFunc(app.createTaskPost))
	router.Handler(http.MethodGet, "/project/create/:companyID", adminInCompanyeis.ThenFunc(app.createProj))
	router.Handler(http.MethodPost, "/project/create/:companyID", adminInCompanyeis.ThenFunc(app.createProjPost))
	router.Handler(http.MethodGet, "/CreateTask/:companyID/:projectName", adminInCompanyeis.ThenFunc(app.createTask))
	router.Handler(http.MethodPost, "/CreateTask/:companyID/:projectName", adminInCompanyeis.ThenFunc(app.createTaskPost))

	OwnerInCompanyeis := protected.Append(app.isOwner)

	router.Handler(http.MethodPost, "/sendInvite/:UserToInviteID/:companyID", OwnerInCompanyeis.ThenFunc(app.sendInvitePost))
	router.Handler(http.MethodGet, "/company/control/:companyID", OwnerInCompanyeis.ThenFunc(app.companyControl))
	router.Handler(http.MethodPost, "/companyControl/:companyID/:userID/:action", OwnerInCompanyeis.ThenFunc(app.companyControlPost))

	siteAdmin := protected.Append(app.isSiteAdmin)

	router.Handler(http.MethodGet, "/admin/usersTable", siteAdmin.ThenFunc(app.usersTable))
	router.Handler(http.MethodGet, "/admin/companyTable", siteAdmin.ThenFunc(app.companyTable))
	router.Handler(http.MethodPost, "/admin/deleteUser/:userID", siteAdmin.ThenFunc(app.adminDeleteUser))
	router.Handler(http.MethodPost, "/admin/deleteCompany/:companyID", siteAdmin.ThenFunc(app.adminDeleteCompany))
	router.Handler(http.MethodGet, "/admin", siteAdmin.ThenFunc(app.adminPanel))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
