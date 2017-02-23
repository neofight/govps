package mock

import "github.com/neofight/govps/tasks"

type Task struct {
	Error     error
	WasCalled bool
}

func (t *Task) Execute(cxt tasks.Context) error {
	t.WasCalled = true
	return t.Error
}
