package stack_test

import (
	"strconv"
	"testing"

	"github.com/neofight/govps/stack"
)

func TestEmptyStackCount(t *testing.T) {
	var s stack.Stack

	if s.Count() != 0 {
		t.Error("Expected empty stack to have a count of 0")
	}
}

func TestEmptyStackPeep(t *testing.T) {
	var s stack.Stack

	if s.Peep() != "" {
		t.Error("Expected peep of empty stack to return an empty string")
	}
}

func TestEmptyStackPop(t *testing.T) {
	var s stack.Stack

	if s.Pop() != "" {
		t.Error("Expected pop of empty stack to return an empty string")
	}
}

func TestStackPushAndCount(t *testing.T) {
	var s stack.Stack

	for i := 0; i < 10; i++ {
		s.Push(strconv.Itoa(i))

		if s.Count() != i+1 {
			t.Error("Count does not match number of items pushed onto the stack")
		}
	}
}

func TestStackPushAndPeep(t *testing.T) {
	var s stack.Stack

	for i := 0; i < 10; i++ {
		s.Push(strconv.Itoa(i))

		if s.Peep() != strconv.Itoa(i) {
			t.Error("Peep does not match the last item pushed onto the stack")
		}
	}
}

func TestStackPushAndPop(t *testing.T) {
	var s stack.Stack

	for i := 0; i < 10; i++ {
		s.Push(strconv.Itoa(i))
	}

	for i := 9; i >= 0; i-- {

		expected := strconv.Itoa(i)
		popped := s.Pop()

		if popped != expected {
			t.Errorf("Expected to pop value %v but saw value %v", expected, popped)
		}
	}
}
