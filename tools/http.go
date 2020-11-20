package tools

import (
	"encoding/json"
	"log"
	"net/http"
)

type e struct {
	Message string `json:"message"`
}

func WriteJsonNotFoundRequest(rw http.ResponseWriter, message string) {
	writeJson(rw, http.StatusNotFound, &e{Message: message})
}

func WriteJsonBadRequest(rw http.ResponseWriter, message string) {
	writeJson(rw, http.StatusBadRequest, &e{Message: message})
}

func WriteJsonCreated(rw http.ResponseWriter, res interface{}) {
	writeJson(rw, http.StatusCreated, res)
}

func WriteJsonOk(rw http.ResponseWriter, res interface{}) {
	writeJson(rw, http.StatusOK, res)
}

func writeJson(rw http.ResponseWriter, status int, res interface{}) {
	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(status)
	if res == nil {
		return
	}

	if err := json.NewEncoder(rw).Encode(res); err != nil {
		log.Printf("Error writing response: %s", err)
	}
}
