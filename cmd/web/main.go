package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"letsgo.bepo1337/internal/models"
	"log"
	"net/http"
	"os"
)

type Application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	snippetModel *models.SnippetModel
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

	app := &Application{
		infoLog:      infoLog,
		errorLog:     errorLog,
		snippetModel: &models.SnippetModel{DB: db},
	}

	infoLog.Printf("Starting server on %s\n", config.addr)
	server := &http.Server{
		Addr:     config.addr,
		Handler:  app.initializeRoutes(config),
		ErrorLog: errorLog}
	err = server.ListenAndServe()
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
