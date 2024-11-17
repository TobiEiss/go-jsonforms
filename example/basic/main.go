package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	gojsonforms "github.com/TobiEiss/go-jsonforms"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	schema   = "testdata/basic/schema.json"
	uiSchema = "testdata/basic/uischema.json"
	data     = "testdata/basic/data.json"
)

func main() {
	schemaData, err := os.ReadFile(schema)
	if err != nil {
		panic(err)
	}

	uiSchemaData, err := os.ReadFile(uiSchema)
	if err != nil {
		panic(err)
	}

	dataD, err := os.ReadFile(data)
	if err != nil {
		panic(err)
	}

	site, err := gojsonforms.New(schemaData, uiSchemaData)
	if err != nil {
		panic(err)
	}

	err = site.BindData(dataD)
	if err != nil {
		panic(err)
	}

	html, err := site.Build()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html)
	})
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		result := gojsonforms.ReadForm(r.Form)
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		fmt.Println(string(jsonData))
	})

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
