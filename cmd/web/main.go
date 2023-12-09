package main

import (
	"log"
	"os"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
}

type server struct {
	host string
	port string
	url  string
}

func main() {
	server := server{
		host: "localhost",
		port: "8080",
		url:  "http://localhost:8080",
	}

	app := &application{
		appName: "hacker News",
		server:  server,
		debug:   true,
		infoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
		errLog:  log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile),
	}

	if err := app.listenAndServe(); err != nil {
		log.Fatal(err)
	}
}
