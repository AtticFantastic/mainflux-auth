package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func addKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	spec, err := readSpec(r)
	if err != nil {
		sErr := err.(*domain.AuthError)
		w.WriteHeader(sErr.Code)
		return
	}

	key, err := services.AddKey(header[1], spec)
	if err != nil {
		sErr := err.(*domain.AuthError)
		w.WriteHeader(sErr.Code)
		return
	}

	rep := fmt.Sprintf(`{"key":"%s"}`, key)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rep))
}

func readSpec(r *http.Request) (domain.KeySpec, error) {
	data := domain.KeySpec{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return data, &domain.AuthError{Code: http.StatusBadRequest}
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return data, &domain.AuthError{Code: http.StatusBadRequest}
	}

	if valid := data.Validate(); !valid {
		return data, &domain.AuthError{Code: http.StatusBadRequest}
	}

	return data, nil
}

func fetchKeys(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	keys, err := services.FetchKeys(header[1])
	if err != nil {
		sErr := err.(*domain.AuthError)
		w.WriteHeader(sErr.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}
