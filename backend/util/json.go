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
	res, err := goreq.Request{
		Method:      "PUT",
		Uri:         url,
		ContentType: "application/json",
		Body:        target,
	}.Do()

	if res != nil {
		res.Body.Close()
	}
	return err
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
