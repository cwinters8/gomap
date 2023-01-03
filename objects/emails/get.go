package emails

import (
	"github.com/cwinters8/gomap/requests"
)

type Get struct {
	requests.Get
	FetchTextBody bool
	FetchHTMLBody bool
}

func (g Get) MarshalJSON() ([]byte, error) {
	m := g.BodyMap()
	m["fetchTextBodyValues"] = g.FetchTextBody
	m["fetchHTMLBodyValues"] = g.FetchHTMLBody
	return requests.MarshalCall(g.ID, "Email/get", m)
}
