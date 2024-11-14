package gojsonforms

type UISchema struct {
	Type     string     `json:"type"`
	Elements UIElements `json:"elements"`
	Label    string     `json:"label"`
}

type UIElement struct {
	Type        string     `json:"type"`
	Scope       string     `json:"scope,omitempty"`
	Text        string     `json:"text,omitempty"`
	Elements    UIElements `json:"elements,omitempty"`
	Suggestions []string   `json:"suggestion,omitempty"`
	Label       string     `json:"label"`
}

type UIElements []UIElement

func (elements UIElements) FindElementWithChild(scope string) UIElements {
	if elements == nil {
		return nil
	}

	for _, child := range elements {
		if child.Scope == scope {
			return elements
		}
		if result := child.Elements.FindElementWithChild(scope); result != nil {
			return result
		}
	}
	return nil
}
