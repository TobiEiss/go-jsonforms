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

var menu = []gojsonforms.MenuItem{
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

		site, err := gojsonforms.New(schemaData, uiSchemaData)
		if err != nil {
			panic(err)
		}

		if dataData != nil {
			err = site.BindData(dataData)
			if err != nil {
				panic(err)
			}
		}

		site.SetMenu(menu)

		html, err := site.Build()
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
