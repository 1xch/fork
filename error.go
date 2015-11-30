package fork

import "fmt"

type forkError struct {
	err  string
	vals []interface{}
}

func (m *forkError) Error() string {
	return fmt.Sprintf("%s", fmt.Sprintf(m.err, m.vals...))
}

func (m *forkError) Out(vals ...interface{}) *forkError {
	m.vals = vals
	return m
}

func Frror(err string) *forkError {
	return &forkError{err: err}
}
