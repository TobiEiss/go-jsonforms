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

var screens = []gojsonforms.Screen{
	{
		Link:  "basic",
		Titel: "Basic Forms",
	},
	{
		Link:  "control",
		Titel: "Control Forms",
	},
	{
		Link:  "array",
		Titel: "Array Forms",
	},
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{screen}", func(w http.ResponseWriter, r *http.Request) {
		screenID := chi.URLParam(r, "screen")
		if screenID == "" {
			screenID = "basic"
		}

		schemaData, err := os.ReadFile(fmt.Sprintf("testdata/%s/schema.json", screenID))
		if err != nil {
			panic(err)
		}

		uiSchemaData, err := os.ReadFile(fmt.Sprintf("testdata/%s/uischema.json", screenID))
		if err != nil {
			panic(err)
		}

		dataData, _ := os.ReadFile(fmt.Sprintf("testdata/%s/data.json", screenID))

		var schema gojsonforms.SchemaJson
		if err := json.Unmarshal(schemaData, &schema); err != nil {
			panic(err)
		}

		var uischema gojsonforms.UIElement
		if err := json.Unmarshal(uiSchemaData, &uischema); err != nil {
			panic(err)
		}

		var html string
		if dataData != nil {
			var data map[string]interface{}
			if err := json.Unmarshal(dataData, &data); err != nil {
				panic(err)
			}
			html, err = gojsonforms.BuildScreenPageWithData(screens, schema, uischema, data)
		} else {
			html, err = gojsonforms.BuildScreenPage(screens, schema, uischema)
		}

		if err != nil {
			panic(err)
		}

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
