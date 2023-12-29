package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/kamogelosekhukhune777/hacker-news-clone/models"
	_ "github.com/lib/pq"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
	view    *jet.Set
	session *scs.SessionManager
	model   models.Models
}

type server struct {
	host string
	port string
	url  string
}

func main() {
	migrate := flag.Bool("migrate", false, "should migrate - drop all tables")
	flag.Parse()
	server := server{
		host: "localhost",
		port: "8080",
		url:  "http://localhost:8080",
	}

	db2, err := openDB("postgres://postgres@...")
	if err != nil {
		log.Fatal(err)
	}
	defer db2.Close()

	//init connection to interact with Postgres
	upper, err := postgresql.New(db2)
	if err != nil {
		log.Fatal(err)
	}
	defer func(upper db.Session) {
		err := upper.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(upper)

	//migrate
	if *migrate {
		fmt.Println("running migration...")
		err := Migrate(upper)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("...Done running migration")
	}

	//init application
	app := &application{
		appName: "hacker News",
		server:  server,
		debug:   true,
		infoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
		errLog:  log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile),
		model:   models.New(upper),
	}

	//init jet template
	if app.debug {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"), jet.InDevelopmentMode())
	} else {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"))
	}

	//init session
	app.session = scs.New()
	app.session.Lifetime = 24 * time.Hour
	app.session.Cookie.Persist = true
	app.session.Cookie.Domain = app.server.host
	app.session.Cookie.SameSite = http.SameSiteStrictMode
	app.session.Store = postgresstore.New(db2)

	if err := app.listenAndServes(); err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db db.Session) error {
	script, err := os.ReadFile(filepath.Join("./migrations/tables.sql"))
	if err != nil {
		return err
	}

	_, err = db.SQL().Exec(script)

	return err
}
