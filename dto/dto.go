package dto

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"io"
	"net/http"
)

type Object interface {
	ToJson(data io.Writer) error
}

type Filter interface {
	ToJson(data io.Writer) error
}

type Update interface {
	ToJson(data io.Writer) error
}

type converter struct {
	Value         interface{}
	ConverterFunc schema.Converter
}

func FilterFromForm(r *http.Request, f Filter, c ...converter) (err error) {
	err = r.ParseForm()
	if err != nil {
		return
	}
	defer r.Body.Close()
	decoder := schema.NewDecoder()
	for _, x := range c {
		decoder.RegisterConverter(x.Value, x.ConverterFunc)
	}
	err = decoder.Decode(f, r.Form)
	return
}

func UpdateFromJson(data io.ReadCloser) (filter Update, err error) {
	err = json.NewDecoder(data).Decode(&filter)
	return
}
