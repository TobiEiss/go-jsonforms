package gojsonforms

type UIElement struct {
	Type        string      `json:"type"`
	Scope       string      `json:"scope,omitempty"`
	Text        string      `json:"text,omitempty"`
	Elements    []UIElement `json:"elements,omitempty"`
	Suggestions []string    `json:"suggestion,omitempty"`
	Label       string      `json:"label"`
}

func (element *UIElement) FindElementWithChild(scope string) *UIElement {
	for i := range element.Elements {
		if element.Elements[i].Scope == scope {
			return element
		}

		if found := element.Elements[i].FindElementWithChild(scope); found != nil {
			return found
		}
	}
	return nil
}
