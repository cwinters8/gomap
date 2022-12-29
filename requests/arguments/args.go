package arguments

type Args interface {
	Query | Set
}

type Filter struct {
	Name string `json:"name"`
}
