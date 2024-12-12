package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time" // New import

	"github.com/AlexPodd/ASAP/internal/models"
	"github.com/alexedwards/scs/mysqlstore" // New import
	"github.com/alexedwards/scs/v2"         // New import
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/justinas/nosurf"
)

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	templateCache     map[string]*template.Template
	users             *models.UserModel
	company           *models.CompanyModel
	usersincompanies  *models.UsersincompaniesModel
	projects          *models.ProjectModel
	tasks             *models.TaskModel
	invites           *models.InviteModel
	formDecoder       *form.Decoder
	sessionManager    *scs.SessionManager
	isDevelopmentMode bool
}

func main() {
	mode := flag.String("mode", "development", "Application mode (development/production)")
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "asapAdmin:password123@/asap?parseTime=true", "MySQL data source name")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	isDevelopmentMode := *mode == "development"
	if isDevelopmentMode {
		infoLog.Printf("Running in development mode")
	} else {
		infoLog.Printf("Running in production mode")
	}

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()
	// And add it to the application dependencies.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app := &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		templateCache:     templateCache,
		users:             &models.UserModel{DB: db},
		company:           &models.CompanyModel{DB: db},
		usersincompanies:  &models.UsersincompaniesModel{DB: db},
		projects:          &models.ProjectModel{DB: db},
		tasks:             &models.TaskModel{DB: db},
		invites:           &models.InviteModel{DB: db},
		formDecoder:       formDecoder,
		sessionManager:    sessionManager,
		isDevelopmentMode: isDevelopmentMode,
	}
	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	// Set the server's TLSConfig field to use the tlsConfig variable we just
	// created.
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
