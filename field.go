package fork

type Field interface {
	New() Field
	Fielder
}

type Fielder interface {
	Name() string
	Get() *Value
	Set(i interface{})
	Processor
}
