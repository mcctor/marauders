package utils

type Collection struct {
	Collection ItemsCollection `json:"collection"`
}

type ItemsCollection struct {
	Version  string            `json:"version"`
	Href     string            `json:"href"`
	Items    []CollectionItem  `json:"items"`
	Links    []CollectionLink  `json:"links"`
	Queries  []CollectionQuery `json:"queries"`
	Template ItemTemplate      `json:"template"`
}

type CollectionItem struct {
	Href  string           `json:"href"`
	Data  []DataField      `json:"data"`
	Links []CollectionLink `json:"links"`
}

type CollectionLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Render string `json:"render"`
}

type CollectionQuery struct {
	Href   string      `json:"href"`
	Rel    string      `json:"rel"`
	Prompt string      `json:"prompt"`
	Data   []DataField `json:"data"`
}

type ItemTemplate struct {
	Data []DataField `json:"data"`
}

type DataField struct {
	Prompt string `json:"prompt"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}
