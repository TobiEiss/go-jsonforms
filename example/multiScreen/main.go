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

var screens = map[string]struct {
	Schema   string
	UISchema string
}{
	"basic": {
		Schema:   "testdata/basic/schema.json",
		UISchema: "testdata/basic/uischema.json",
	},
	"control": {
		Schema:   "testdata/control/schema.json",
		UISchema: "testdata/control/uischema.json",
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

		schemaData, err := os.ReadFile(screens[screenID].Schema)
		if err != nil {
			panic(err)
		}

		uiSchemaData, err := os.ReadFile(screens[screenID].UISchema)
		if err != nil {
			panic(err)
		}

		var schema gojsonforms.SchemaJson
		if err := json.Unmarshal(schemaData, &schema); err != nil {
			panic(err)
		}

		var uischema gojsonforms.UIElement
		if err := json.Unmarshal(uiSchemaData, &uischema); err != nil {
			panic(err)
		}

		screens := []gojsonforms.Screen{
			{
				Titel: "Basic Form",
				Link:  "basic",
			},
			{
				Titel: "Control Form",
				Link:  "control",
			},
		}
		html, err := gojsonforms.BuildScreenPage(screens, schema, uischema)
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
