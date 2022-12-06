package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndPoints(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "/login", LoginPath)
	assert.Equal(t, "/register", RegisterPath)
	assert.Equal(t, "/activate", ActivateAccountPath)
	assert.Equal(t, "/projects", ProjectsPath)
	assert.Equal(t, "/projects/{id}", ProjectPath)
	assert.Equal(t, "/projects/{projectId}/tasks", TasksPath)
	assert.Equal(t, "/projects/{projectId}/tasks/{taskId}", TaskPath)
}
