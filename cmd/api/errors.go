package main

import (
	"log"
	"net/http"
)

func (app *application) permissionDenied(w http.ResponseWriter, r *http.Request) {
	log.Printf("forbidden request: %s, path: %s", r.Method, r.URL.Path)
	WriteJSONError(w, http.StatusForbidden, "permission denied")
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("unauthorized error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	WriteJSONError(w, http.StatusInternalServerError, "server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	WriteJSONError(w, http.StatusNotFound, "not found")
}
