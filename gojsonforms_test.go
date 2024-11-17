package gojsonforms_test

import (
	"reflect"
	"testing"

	gabs "github.com/Jeffail/gabs/v2"
	gojsonforms "github.com/TobiEiss/go-jsonforms"
)

func TestIteration(t *testing.T) {
	tests := []struct {
		testStep string
		schema   string
		uiSchema string
		data     string
		expected string
	}{
		{
			testStep: "simple object",
			schema: `{
					"type": "object",
					"properties": {
						"name": {
							"type":        "string",
							"minlength":   3,
							"description": "enter name"
						}
					}
				}`,
			uiSchema: `{
					"type": "verticalLayout",
					"elements": [
						{
							"type":  "Control",
							"scope": "#/properties/name"
						}
					]
				}`,
			data: `{
					"name": "John Doe",
					"vegetarian": false,
					"birthDate": "1985-06-02",
					"personalData": {
						"age": 34
					},
					"postalCode": "12345"
				}`,
			expected: `{
					"type": "verticalLayout",
					"elements": [
						{
							"type":  "Control",
							"scope": "#/properties/name",
							"schema": {
								"type":        "string",
								"minlength":   3,
								"description": "enter name"
							},
							"data": "John Doe"
						}
					]
				}`,
		},
		{
			testStep: "with enum object",
			schema: `{
					"type": "object",
					"properties": {
						"country": {
							"enum":        ["DE", "IT", "JP"],
							"description": "enter country"
						}
					}
				}`,
			uiSchema: `{
					"type": "verticalLayout",
					"elements": [
						{
							"type":  "Control",
							"scope": "#/properties/country"
						}
					]
				}`,
			expected: `{
					"type": "verticalLayout",
					"elements": [
						{
							"type":  "Control",
							"scope": "#/properties/country",
							"schema": {
								"enum":        ["DE", "IT", "JP"],
								"description": "enter country"
							}
						}
					]
				}`,
		},
		{
			testStep: "array",
			schema: `{
				"properties": {
					"comments": {
						"type": "array",
						"title": "Comments",
						"items": {
							"type": "object",
							"properties": {
								"message": {
        							"type": "string"
        						},
    							"name": {
            						"type": "string"
          						}
        					}
    					}
    				}
  				}				
  			}`,
			uiSchema: `{
  				"type": "VerticalLayout",
  				"elements": [
    				{
      					"type": "Control",
      					"scope": "#/properties/comments",
      					"options": {
        					"elementLabelProp": "name",
        					"detail": {
          						"type": "HorizontalLayout",
          						"elements": [
						            {
						              	"type": "Control",
						              	"scope": "#/properties/comments/items/properties/message"
						            },
						            {
							            "type": "Control",
						              	"scope": "#/properties/comments/items/properties/name"
						            }
						        ]
						    }
						}
					}
				]
  			}`,
			data: `{
  				"comments": [
					{
				    	"name": "John Doe",
				      	"message": "This is an example message"
				    },
				    {
				      	"name": "Max Mustermann",
				      	"message": "Another message"
				    }
				]
			}`,
			expected: `{
  				"type": "VerticalLayout",
  				"elements": [
    				{
      					"type": "Control",
      					"scope": "#/properties/comments",
      					"schema": {
      						"type": "array",
      						"title": "Comments"	
      					},
      					"options": {
        					"elementLabelProp": "name",
        					"detail": [
        						{
	        						"type": "HorizontalLayout",
	        						"elements": [
							            {
							              	"type": "Control",
							              	"scope": "#/properties/comments/0/properties/message",
							              	"schema": {
							              		"type": "string",
							              		"col": " column col-6"
							              	},
							              	"data": "This is an example message"
							            },
							            {
								            "type": "Control",
							              	"scope": "#/properties/comments/0/properties/name",
							              	"schema": {
							              		"type": "string",
							              		"col": " column col-6"							              	},
							              	"data": "John Doe"
							            }
									]
	        					},
        						{
	        						"type": "HorizontalLayout",
	        						"elements": [
							            {
							              	"type": "Control",
							              	"scope": "#/properties/comments/1/properties/message",
							              	"schema": {
							              		"type": "string",
							              		"col": " column col-6"							              	},
							              	"data": "Another message"
							            },
							            {
								            "type": "Control",
							              	"scope": "#/properties/comments/1/properties/name",
							              	"schema": {
							              		"type": "string",
							              		"col": " column col-6"							              	},
							              	"data": "Max Mustermann"
							            }
									]
	        					}
	        					
        					]
						}
					}
				]
  			}`,
		},
	}

	// gojsonforms.Iterate()

	for _, test := range tests {
		t.Run(test.testStep, func(t *testing.T) {
			site, err := gojsonforms.New([]byte(test.schema), []byte(test.uiSchema))
			if err != nil {
				t.Error(err)
				return
			}

			if test.data != "" {
				err = site.BindData([]byte(test.data))
				if err != nil {
					t.Error(err)
					return
				}
			}

			toJsonString := func(input []byte) string {
				obj, err := gabs.ParseJSON(input)
				if err != nil {
					t.Error("no valid json", err)
				}
				return obj.String()
			}

			uiSchemaString := toJsonString(site.UISchema())
			expectedString := toJsonString([]byte(test.expected))
			if !reflect.DeepEqual(uiSchemaString, expectedString) {
				t.Errorf("not equal:\n%s\n%s", uiSchemaString, expectedString)
			}
		})
	}
}
