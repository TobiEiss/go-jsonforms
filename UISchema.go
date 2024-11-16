package gojsonforms

type UIElement struct {
	Type        string      `json:"type"`
	Scope       string      `json:"scope,omitempty"`
	Text        string      `json:"text,omitempty"`
	Elements    []UIElement `json:"elements,omitempty"`
	Suggestions []string    `json:"suggestion,omitempty"`
	Label       string      `json:"label"`
	Options     Options     `json:"options"`
}

type Options struct {
	ElementLabelProp string `json:"elementLabelProp"`
	Detail           Detail `json:"detail"`
}

type Detail struct {
	Type     string      `json:"type"`
	Elements []UIElement `json:"elements,omitempty"`
}

func (element *UIElement) FindParentWithScope(scope string) (string, []UIElement) {
	for i := range element.Elements {
		if element.Elements[i].Scope == scope {
			return element.Type, element.Elements
		}

		if found, elements := element.Elements[i].FindParentWithScope(scope); found != "" {
			return found, elements
		}
	}

	for i := range element.Options.Detail.Elements {
		if element.Options.Detail.Elements[i].Scope == scope {
			return element.Options.Detail.Type, element.Options.Detail.Elements
		}

		if found, elements := element.Options.Detail.Elements[i].FindParentWithScope(scope); found != "" {
			return found, elements
		}
	}

	return "", []UIElement{}
}

func (element *UIElement) FindWithScope(scope string) *UIElement {
	if element.Scope == scope {
		return element
	}

	for i := range element.Elements {
		if found := element.Elements[i].FindWithScope(scope); found != nil {
			return found
		}
	}

	for i := range element.Options.Detail.Elements {
		if found := element.Options.Detail.Elements[i].FindWithScope(scope); found != nil {
			return found
		}
	}

	return nil
}
