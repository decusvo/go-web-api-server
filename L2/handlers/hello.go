package handlers // defines the fact that this file will be part of a package called 'handlers'.

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello { // a function that takes in a log.logger object and returns a Hello handler as a reference.
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	h.l.Println("Hello World!")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "ERROR", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "Hello %s", d)

}
