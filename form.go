package fork

import "net/http"

type Form interface {
	New() Form
	Former
}

type Former interface {
	Fields(...Field) []Field
	Process(*http.Request)
	Render() string
	Valid() bool
}
