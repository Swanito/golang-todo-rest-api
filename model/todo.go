package model

import (
	"time"
)

type Projects struct {
	Projects []*ProjectResponse `json:"projects"`
}

type Tasks struct {
	Tasks []*TaskResponse `json:"tasks"`
}

type Task struct {
	Id          int       `json:"id,omitempty"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	ProjectId   int       `json:"projectId,omitempty"`
	UserId      int       `json:"userId,omitempty"`
}

type TaskResponse struct {
	Id          int       `json:"id,omitempty"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	UserId      int       `json:"createdBy"`
}

type ProjectResponse struct {
	Title       string          `json:"name,omitempty"`
	Completed   bool            `json:"completed"`
	CompletedAt time.Time       `json:"completedAt,omitempty"`
	Id          int             `json:"id,omitempty"`
	UserId      int             `json:"createdBy"`
	Tasks       []*TaskResponse `json:"tasks,omitempty"`
}

type CreateTaskRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type CreateUpdateProjectRequest struct {
	Title string `json:"title,omitempty"`
}

type TodoRepository interface {
	GetProjects(userId int) ([]*ProjectResponse, error)
	GetProject(userId int, projectId int) (*ProjectResponse, error)
	CreateProject(userId int, title string) (*ProjectResponse, error)
	DeleteProject(userId, projectId int) error
	UpdateProject(userId, projectId int, title string) (*ProjectResponse, error)
	CompleteProject(userId, projectId int) (*ProjectResponse, error)
	CreateTask(userId, projectId int, title, description string) (*TaskResponse, error)
	GetTasks(userId, projectId int) ([]*TaskResponse, error)
	CompleteTask(userId, projectId, taskId int) (*TaskResponse, error)
}
