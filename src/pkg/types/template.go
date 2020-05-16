package types

// TemplateList is the view model for template list page
type TemplateList struct {
	Base
	Templates []Template
}

// Template represents a workflow template
type Template struct {
	ID          string
	Name        string
	Data        string
	LastUpdated string
}
