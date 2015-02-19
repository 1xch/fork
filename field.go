package fork

import "net/http"

type Field interface {
	New() Field
	Fielder
}

type Fielder interface {
	Name(...string) string
	Get() *Value
	Set(r *http.Request)
	Processor
}
