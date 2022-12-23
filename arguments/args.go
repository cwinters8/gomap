package arguments

type Args interface {
	Query
}

type Filter struct {
	Name string `json:"name"`
}
