package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/franela/goreq"
)

func GetJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func PutJson(url string, target interface{}) error {
	r, err := goreq.Request{
		Method:      "PUT",
		Uri:         url,
		ContentType: "application/json",
		Body:        target,
	}.Do()
	defer r.Body.Close()

	return err
}

func PostJson(url string, body interface{}, response interface{}) error {
	r, err := goreq.Request{
		Method:      "POST",
		Uri:         url,
		ContentType: "application/json",
		Body:        body,
	}.Do()
	defer r.Body.Close()

	if err != nil {
		return err
	}

	return json.NewDecoder(r.Body).Decode(response)
}

func PrettyPrintJson(target interface{}) {
	b, err := json.Marshal(target)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, " ", " ")
	out.WriteTo(os.Stdout)
}
