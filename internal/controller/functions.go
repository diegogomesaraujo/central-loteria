package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/diegogomesaraujo/central-loteria/pkg/exception"
	"github.com/diegogomesaraujo/central-loteria/pkg/repository"
)

func readBodyFromJSON(w http.ResponseWriter, r *http.Request, entity interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		handleExceptionError(w, "Content-Type not supported", http.StatusUnsupportedMediaType)
		return fmt.Errorf("Content-Type not supported")
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error when read request body: %v\n", err)
		handleExceptionError(w, "Erro when try read body", http.StatusInternalServerError)
		return fmt.Errorf("Erro when try read body")
	}

	if !json.Valid(body) {
		handleExceptionError(w, "Bad Request", http.StatusBadRequest)
		return fmt.Errorf("Bad Request")
	}

	json.Unmarshal(body, entity)

	return nil
}

func handleExceptionError(w http.ResponseWriter, message string, httpCode int) {
	ex := exception.Exception{
		Message: message,
		Code:    httpCode,
	}
	exception.HandleError(w, ex)
}

func firestoreConnect(w http.ResponseWriter) (repository.Firestore, error) {
	firestore := repository.Firestore{}
	err := firestore.Connect()

	if err != nil {
		log.Printf("Error when connect to firestore: %v\n", err)
		handleExceptionError(w, "Connection error", http.StatusInternalServerError)
		return firestore, errors.New("Connection error")
	}

	return firestore, nil
}