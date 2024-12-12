package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.

		//w.Header().Set("Content-Security-Policy",
		//	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com; script-src 'self'")

		nonce := generateNonce()
		w.Header().Set("Content-Security-Policy", fmt.Sprintf("script-src 'self' 'nonce-%s'; style-src 'self' fonts.googleapis.com", nonce))
		ctx := context.WithValue(r.Context(), "nonce", nonce)
		r = r.WithContext(ctx)

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")
		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func (app *application) isMember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idUser := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if idUser == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		params := httprouter.ParamsFromContext(r.Context())

		companyID, err := strconv.Atoi(params.ByName("companyID"))
		if err != nil || companyID < 1 {
			app.notFound(w)
			return
		}

		isMember, err := app.usersincompanies.IsUserInCompany(idUser, companyID)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !isMember {
			app.sessionManager.Put(r.Context(), "flash", "You must be in company to perform this action!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Передача управления следующему обработчику
		next.ServeHTTP(w, r)
	})
}
func (app *application) isAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idUser := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if idUser == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		params := httprouter.ParamsFromContext(r.Context())
		companyID, err := strconv.Atoi(params.ByName("companyID"))
		if err != nil || companyID < 1 {
			app.errorLog.Print(err)
			app.notFound(w)
			return
		}

		isAdmin, err := app.usersincompanies.IsUserAAdminOrOwner(idUser, companyID)
		if err != nil {
			app.errorLog.Print(err)
			app.serverError(w, err)
			return
		}

		if !isAdmin {
			app.sessionManager.Put(r.Context(), "flash", "You must be admin/owner in company to perform this action!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Передача управления следующему обработчику
		next.ServeHTTP(w, r)
	})
}
func (app *application) isOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idUser := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if idUser == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		params := httprouter.ParamsFromContext(r.Context())
		//asdasdasdasd
		companyID, err := strconv.Atoi(params.ByName("companyID"))
		if err != nil || companyID < 1 {
			app.notFound(w)
			return
		}

		isOwner, err := app.usersincompanies.IsUserAOwner(idUser, companyID)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !isOwner {
			app.sessionManager.Put(r.Context(), "flash", "You must be in owner of company to perform this action!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Передача управления следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func (app *application) isSiteAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idUser := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

		isSiteAdmin, err := app.users.IsUserASiteAdmin(idUser)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !isSiteAdmin {
			app.sessionManager.Put(r.Context(), "flash", "You must be in site admin to perform this action!")
			app.sessionManager.Put(r.Context(), "flashtype", "error")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Передача управления следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func generateNonce() string {
	nonce := make([]byte, 16) // 16 байт для nonce
	_, err := rand.Read(nonce)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(nonce)
}
