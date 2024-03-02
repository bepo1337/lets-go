package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"letsgo.bepo1337/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	Cost = 10
)

type Application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippetModel   models.SnippetModelInterface
	userModel      models.UserModelInterface
	templates      map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

type Config struct {
	addr      string
	staticDir string
}

func setupConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.addr, "addr", "localhost:4000", "HTTP port")
	flag.StringVar(&cfg.staticDir, "dir", "./ui/static/", "Directory for serving static files")
	return &cfg
}

func main() {
	config := setupConfig()
	dataSourceString := flag.String("dsn", "root:admin@/snippetbox?parseTime=true", "MySQL data source string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDbConnectionPool(dataSourceString)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templates, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	sessionManager := scs.New() //sets default values and calls constructor
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &Application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		snippetModel:   &models.SnippetModel{DB: db},
		userModel:      &models.UserModel{DB: db},
		templates:      templates,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	server := &http.Server{
		Addr:         config.addr,
		Handler:      app.initializeRoutes(),
		ErrorLog:     errorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s\n", config.addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)

}

func openDbConnectionPool(dataSourceString *string) (*sql.DB, error) {
	db, err := sql.Open("mysql", *dataSourceString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
