package tasks_test

import (
	"errors"
	"testing"

	"github.com/neofight/govps/mock"
	"github.com/neofight/govps/tasks"
)

func TestPipelineSuccess(t *testing.T) {

	server := mock.NewServer()

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	mockTasks := make([]tasks.Task, 10)

	for i := range mockTasks {
		mockTasks[i] = new(mock.Task)
	}

	err := tasks.ExecutePipeline(cxt, mockTasks)

	if err != nil {
		t.Error("Expected the tasks to be run without error but they were not")
	}
}

var taskError = errors.New("Task failed")

func TestPipelineFailure(t *testing.T) {

	server := mock.NewServer()

	cxt := tasks.Context{VPS: server, Domain: "test.com"}

	mockTasks := make([]tasks.Task, 10)

	for i := range mockTasks {
		mockTasks[i] = new(mock.Task)
	}

	mockTasks[5].(*mock.Task).Error = taskError

	err := tasks.ExecutePipeline(cxt, mockTasks)

	if err != taskError {
		t.Error("Expected the pipleline to return the error of the failed task but it did not")
	}

	for i := 0; i < 10; i++ {
		if i <= 5 {
			if !mockTasks[i].(*mock.Task).WasCalled {
				t.Error("Expected all tasks up to and including the failed task to have run")
			}
			continue
		}

		if mockTasks[i].(*mock.Task).WasCalled {
			t.Error("Expected all tasks after the failed task to have been skipped")
		}
	}
}
