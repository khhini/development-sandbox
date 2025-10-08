package domain

import "testing"

func TestNewTask(t *testing.T) {
	testCase := struct {
		title       string
		description string
	}{
		title:       "New Test Task",
		description: "Creating new task for testing",
	}

	task := NewTask(testCase.title, testCase.description)

	if testCase.title != task.Title {
		t.Errorf("task.Title = %q, want %q", task.Title, testCase.title)
	}

	if task.ID == "" {
		t.Errorf("task.ID = %q, task.id must not empty", task.ID)
	}
}
