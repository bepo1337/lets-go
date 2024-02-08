package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
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
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &Application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s\n", config.addr)
	server := &http.Server{
		Addr:     config.addr,
		Handler:  app.initializeRoutes(config),
		ErrorLog: errorLog}
	err := server.ListenAndServe()
	errorLog.Fatal(err)

}
