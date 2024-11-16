package gojsonforms

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	gabs "github.com/Jeffail/gabs/v2"
)

//go:embed html/*
var resources embed.FS

type SchemaJson map[string]interface{}
type Screen struct {
	Titel string
	Link  string
}

func BuildSinglePage(schema SchemaJson, uiSchema UIElement) (string, error) {
	return buildTemplate([]Screen{}, schema, uiSchema, nil)
}

func BuildScreenPage(screens []Screen, schema SchemaJson, uiSchema UIElement) (string, error) {
	return buildTemplate(screens, schema, uiSchema, nil)
}

func BuildScreenPageWithData(screens []Screen, schema SchemaJson, uiSchema UIElement, data map[string]interface{}) (string, error) {
	return buildTemplate(screens, schema, uiSchema, data)
}

func buildTemplate(screens []Screen, schema SchemaJson, uiSchema UIElement, data map[string]interface{}) (string, error) {
	var builder strings.Builder
	var err error

	find := func(scope string) SchemaJson {
		re := regexp.MustCompile(`\d`)
		scope = re.ReplaceAllString(scope, "items")
		pathSplits := strings.Split(scope, "/")
		element, findErr := findElement(pathSplits[1:], schema)
		if findErr != nil {
			slog.Error(fmt.Sprintf("%s in %s", findErr.Error(), scope))
		}

		element["Label"] = pathSplits[len(pathSplits)-1]
		element["Scope"] = scope

		return element
	}

	colWidthClass := func(scope string) string {
		parentsType, elements := uiSchema.FindParentWithScope(scope)
		if parentsType == "" {
			slog.Error(fmt.Sprintf("can't find %s", scope))
		}

		if parentsType != "HorizontalLayout" {
			return ""
		}

		return fmt.Sprintf(" column col-%d", 12/len(elements))
	}

	findData := func(scope string) template.HTMLAttr {
		scope = strings.TrimPrefix(scope, "#/")
		scope = strings.ReplaceAll(scope, "properties/", "")
		// scope = strings.ReplaceAll(scope, "items/", "1/")
		path := strings.Split(scope, "/")

		jsonParsed := gabs.Wrap(data)
		value := jsonParsed.Search(path...).Data()

		return template.HTMLAttr(fmt.Sprintf("value=\"%s\"", value))
	}

	generateItems := func(scope string) []UIElement {
		element := uiSchema.FindWithScope(scope)
		if element == nil {
			return []UIElement{}
		}

		items := element.Options.Detail.Elements
		for i := range items {
			replace := fmt.Sprintf("/%d", i)
			items[i].Scope = strings.Replace(items[i].Scope, "/items", replace, 1)
		}
		fmt.Printf("items: %v\n", items)
		return items
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"find":          find,
		"colWidthClass": colWidthClass,
		"findData":      findData,
		"generateItems": generateItems,
	}).ParseFS(resources, "html/*")
	if err != nil {
		return builder.String(), err
	}

	err = tmpl.ExecuteTemplate(&builder, "index.html", map[string]interface{}{
		"Screens":  screens,
		"UISchema": uiSchema,
		"Schema":   schema,
	})

	return builder.String(), err
}

func ReadForm(form url.Values) string {
	jsonObj := gabs.New()

	for key, value := range form {
		path := strings.TrimPrefix(key, "#/")
		path = strings.ReplaceAll(key, "properties/", "")
		keys := strings.Split(path, "/")

		val := value[0]

		if numVal, err := strconv.Atoi(val); err == nil {
			jsonObj.Set(numVal, keys...)
		} else {
			jsonObj.Set(val, keys...)
		}
	}

	return jsonObj.String()
}

func findElement(path []string, s map[string]interface{}) (SchemaJson, error) {
	if len(path) == 0 {
		return s, nil
	}

	if _, ok := s[path[0]]; !ok {
		return nil, errors.New(fmt.Sprintf("Can't find %s", s[path[0]]))
	}
	return findElement(path[1:], s[path[0]].(map[string]interface{}))
}
