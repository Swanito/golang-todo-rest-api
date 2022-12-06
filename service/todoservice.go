package service

import (
	"errors"
	"go-playground/model"
	"net/http"
)

func GetProjects(userId int, repository model.TodoRepository) ([]*model.ProjectResponse, *model.CustomError) {
	projects, err := repository.GetProjects(userId)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return projects, nil
}

func GetProject(userId int, projectId int, repository model.TodoRepository) (*model.ProjectResponse, *model.CustomError) {
	project, err := repository.GetProject(userId, projectId)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	if project.Title == "" {
		return nil, &model.CustomError{Error: errors.New("Project not found"), StatusCode: http.StatusNotFound}
	}
	return project, nil
}

func CreateProject(userId int, title string, todoR model.TodoRepository) (*model.ProjectResponse, *model.CustomError) {
	project, err := todoR.CreateProject(userId, title)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return project, nil
}

func DeleteProject(userId int, projectId int, repository model.TodoRepository) *model.CustomError {
	err := repository.DeleteProject(userId, projectId)
	if err != nil {
		return &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return nil
}

func UpdateProject(userId, projectId int, newTitle string, repository model.TodoRepository) (*model.ProjectResponse, *model.CustomError) {
	project, err := repository.UpdateProject(userId, projectId, newTitle)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return project, nil
}

func CompleteProject(userId, projectId int, repository model.TodoRepository) (*model.ProjectResponse, *model.CustomError) {
	project, err := repository.CompleteProject(userId, projectId)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return project, nil
}

func CreateTask(userId, projectId int, title, description string, repository model.TodoRepository) (*model.TaskResponse, *model.CustomError) {
	task, err := repository.CreateTask(userId, projectId, title, description)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return task, nil
}

func GetTasksForProject(userId, projectId int, repository model.TodoRepository) (*model.Tasks, *model.CustomError) {
	tasks, err := repository.GetTasks(userId, projectId)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return &model.Tasks{Tasks: tasks}, nil
}

func CompleteTask(userId, projectId, taskId int, repository model.TodoRepository) (*model.TaskResponse, *model.CustomError) {
	task, err := repository.CompleteTask(userId, projectId, taskId)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return task, nil
}
