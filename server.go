package main

import (
	"fmt"
	"net/http"
	"time"
)

// starts the server
func (a *application) listenAndServes() error {
	host := fmt.Sprintf("%s:%s", a.server.host, a.server.port)
	serv := http.Server{
		Handler:     a.routes(),
		Addr:        host,
		ReadTimeout: 300 * time.Second,
	}

	a.infoLog.Printf("server listening on :%s\n", host)

	return serv.ListenAndServe()
}
