package backend

import (
	"net/http"
	"strconv"

	"github.com/mindcastio/go-json-rest/rest"

	"github.com/mindcastio/mindcastio/backend/jsonapi"
	"github.com/mindcastio/mindcastio/backend/util"
)

type (
	JsonApiError struct {
		Errors *[]Error `json:"errors"`
	}

	Resource struct {
		Id         string                 `json:"id"`
		Type       string                 `json:"type"`
		Attributes map[string]interface{} `json:"attributes,omitempty"`
	}

	Error struct {
		Id     string `json:"id"`
		Status string `json:"status"`
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}
)

func Response(w rest.ResponseWriter, model interface{}) error {
	w.WriteHeader(http.StatusOK)
	err := w.WriteJson(model)

	return err
}

func StatusResponse(w rest.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func JsonApiResponse(w rest.ResponseWriter, model interface{}) error {
	payload, err := jsonapi.MarshalOne(model)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	err = w.WriteJson(payload)

	return err
}

func JsonApiArrayResponse(w rest.ResponseWriter, models []interface{}) error {
	payload, err := jsonapi.MarshalMany(models)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	err = w.WriteJson(payload)

	return err
}

func JsonApiErrorResponse(w rest.ResponseWriter, code string, message string, err error) {
	var msg string = message
	if err != nil {
		msg = err.Error()
	}

	uuid, _ := util.UUID()
	errors := make([]Error, 1)
	errors[0] = Error{
		uuid,
		strconv.Itoa(http.StatusBadRequest),
		code,
		message,
		msg,
	}

	w.WriteHeader(http.StatusBadRequest)
	ee := w.WriteJson(JsonApiError{&errors})
	if ee != nil {
		panic(ee)
	}
}
