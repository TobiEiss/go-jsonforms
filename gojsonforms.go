package gojsonforms

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	gabs "github.com/Jeffail/gabs/v2"
)

//go:embed html/*
var resources embed.FS

type Form struct {
	schema   *gabs.Container
	uiSchema *gabs.Container
	data     *gabs.Container
	menu     []MenuItem
}

type MenuItem struct {
	Link  string
	Titel string
}

func New(schema []byte, uiSchema []byte) (Form, error) {
	var form Form
	schemaParsed, err := gabs.ParseJSON(schema)
	if err != nil {
		return form, err
	}

	uiSchemaParsed, err := gabs.ParseJSON(uiSchema)
	if err != nil {
		return form, err
	}

	form.schema = schemaParsed
	form.uiSchema = uiSchemaParsed
	err = form.setup()

	return form, err
}

func (form *Form) setup() error {
	var err error

	// add schema-information as schema to every control
	iterateObj(form.uiSchema, "type", "Control", func(c *gabs.Container) {
		scope, ok := c.Path("scope").Data().(string)
		if !ok {
			err = errors.New(fmt.Sprintf("in %v is no scope", c))
		}

		for k, v := range form.schema.Path(gabsPath(scope, true)).ChildrenMap() {
			// "simple" (not nested) object
			if len(v.Children()) == 0 {
				c.SetP(v, fmt.Sprintf("schema.%s", k))
			}
			// arrays (for e.g. items)
			if _, err := v.ArrayCount(); err == nil {
				c.SetP(v, fmt.Sprintf("schema.%s", k))
			}
		}
	})

	// add HTML-col-tag
	iterateObj(form.uiSchema, "type", "HorizontalLayout", func(c *gabs.Container) {
		arrayCount, err := c.ArrayCountP("elements")
		if err != nil {
			fmt.Println(err)
			return
		}

		tag := fmt.Sprintf(" column col-%d", 12/arrayCount)
		for i := range arrayCount {
			c.SetP(tag, fmt.Sprintf("elements.%d.schema.col", i))
		}
	})

	return err
}

func (form *Form) BindData(data []byte) error {
	var err error

	dataParsed, err := gabs.ParseJSON(data)
	form.data = dataParsed

	// build multiple items for arrays
	arrayObj, _ := gabs.New().Set("array")
	iterateObj(form.uiSchema, "schema.type", arrayObj, func(c *gabs.Container) {
		scope, ok := c.Path("scope").Data().(string)
		if !ok {
			err = errors.New(fmt.Sprintf("can't find scope for %v", c))
		}

		// how many data are there
		arrayCount, err := form.data.ArrayCountP(gabsPath(scope, false))
		if err != nil {
			return
		}

		// create new array
		origin := c.Path("options.detail").String()
		for i := range arrayCount {
			copy := strings.ReplaceAll(origin, "items/", fmt.Sprintf("%d/", i))
			newDetails, _ := gabs.ParseJSON([]byte(copy))
			c.ArrayAppendP(newDetails, "options.detail")
		}
		c.ArrayRemoveP(0, "options.detail")
	})

	// I don't know why....
	form.uiSchema, _ = gabs.ParseJSON([]byte(form.uiSchema.String()))

	// add data to every control
	iterateObj(form.uiSchema, "type", "Control", func(c *gabs.Container) {
		// ignore array-controls
		schemaType := c.Path("schema.type").Data()
		if reflect.DeepEqual(schemaType, "array") {
			return
		}

		scope, ok := c.Path("scope").Data().(string)
		if !ok {
			err = errors.New(fmt.Sprintf("in %v is no scope", c))
		}

		c.SetP(form.data.Path(gabsPath(scope, false)).Data(), "data")
	})
	return err
}

func (form *Form) SetMenu(menu []MenuItem) {
	form.menu = menu
}

func (form *Form) Build() (string, error) {
	var builder strings.Builder
	var err error

	tmpl, err := template.New("").ParseFS(resources, "html/*")
	if err != nil {
		return builder.String(), err
	}

	var uischema map[string]interface{}
	if err := json.Unmarshal(form.uiSchema.Bytes(), &uischema); err != nil {
		panic(err)
	}

	err = tmpl.ExecuteTemplate(&builder, "index.html", map[string]interface{}{
		"UISchema": uischema,
		"Menu":     form.menu,
	})

	return builder.String(), err
}

func ReadForm(urlForm url.Values) string {
	jsonObj := gabs.New()

	for key, value := range urlForm {
		path := gabsPath(key, false)

		val := value[0]
		if numVal, err := strconv.Atoi(val); err == nil {
			jsonObj.SetP(numVal, path)
		} else {
			jsonObj.SetP(val, path)
		}
	}

	return jsonObj.String()
}

func (form *Form) UISchema() []byte {
	return form.uiSchema.Bytes()
}

// iterateObj searches for a key and value. If value is empty, it looks only for the key
func iterateObj(container *gabs.Container, key string, value any, operate func(c *gabs.Container)) {
	val := container.Path(key).Data()
	if val != nil && (value == nil || reflect.DeepEqual(val, value)) {
		operate(container)
	}

	for _, child := range container.Children() {
		iterateObj(child, key, value, operate)
	}
}

func gabsPath(scope string, withProperties bool) string {
	scope = strings.Trim(scope, "#/")
	if !withProperties {
		scope = strings.ReplaceAll(scope, "properties/", "")
	}
	path := strings.ReplaceAll(scope, "/", ".")
	return path
}
