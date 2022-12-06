package handlers

import (
	"encoding/json"
	"go-playground/model"
	"go-playground/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *BaseHandler) CreateProjectHandler(w http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	createProjectRequest := model.CreateUpdateProjectRequest{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&createProjectRequest); err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	project, err := service.CreateProject(userId, createProjectRequest.Title, h.todoRepository)
	if err != nil {
		h.SendError(w, err.StatusCode, err.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, project)
}

func (h *BaseHandler) GetProjectHandler(w http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	params := mux.Vars(req)
	projectId, err := strconv.Atoi(params["id"])
	if err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	project, projectErr := service.GetProject(userId, projectId, h.todoRepository)
	if projectErr != nil {
		h.SendError(w, projectErr.StatusCode, projectErr.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, project)
}

func (h *BaseHandler) GetProjectsHandler(w http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	projects, _ := service.GetProjects(userId, h.todoRepository)
	h.Send(w, http.StatusOK, model.Projects{Projects: projects})
}

func (h *BaseHandler) DeleteProjectHandler(w http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	params := mux.Vars(req)
	projectId, err := strconv.Atoi(params["id"])
	if err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	projectErr := service.DeleteProject(userId, projectId, h.todoRepository)
	if projectErr != nil {
		h.SendError(w, projectErr.StatusCode, projectErr.Error.Error())
		return
	}
	h.Send(w, http.StatusNoContent, nil)
}

func (h *BaseHandler) UpdateProjectHandler(w http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	params := mux.Vars(req)
	projectId, err := strconv.Atoi(params["id"])
	if err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	createProjectRequest := model.CreateUpdateProjectRequest{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&createProjectRequest); err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	project, updateErr := service.UpdateProject(userId, projectId, createProjectRequest.Title, h.todoRepository)
	if updateErr != nil {
		h.SendError(w, updateErr.StatusCode, updateErr.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, project)
}

func (h *BaseHandler) CompleteProjectHandler(w http.ResponseWriter, req *http.Request) {
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	params := mux.Vars(req)
	projectId, err := strconv.Atoi(params["id"])
	if err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	project, completeErr := service.CompleteProject(userId, projectId, h.todoRepository)
	if completeErr != nil {
		h.SendError(w, completeErr.StatusCode, completeErr.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, project)
}

func (h *BaseHandler) CreateTaskHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	projectId, err := strconv.Atoi(params["projectId"])
	if err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	createTaskRequest := model.CreateTaskRequest{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&createTaskRequest); err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	project, projectErr := service.GetProject(userId, projectId, h.todoRepository)
	if projectErr != nil {
		h.SendError(w, projectErr.StatusCode, projectErr.Error.Error())
		return
	}
	if project.Completed {
		h.SendError(w, http.StatusBadRequest, "Cannot add new tasks on a completed project")
		return
	}
	task, createErr := service.CreateTask(userId, projectId, createTaskRequest.Title, createTaskRequest.Description, h.todoRepository)
	if createErr != nil {
		h.SendError(w, createErr.StatusCode, createErr.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, task)
}

func (h *BaseHandler) GetTasksHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	projectId, err := strconv.Atoi(params["projectId"])
	if err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	tasks, createErr := service.GetTasksForProject(userId, projectId, h.todoRepository)
	if createErr != nil {
		h.SendError(w, createErr.StatusCode, createErr.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, tasks)
}

func (h *BaseHandler) CompleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	userId, _ := strconv.Atoi(req.Header["Userid"][0])
	projectId, err := strconv.Atoi(params["projectId"])
	taskId, err := strconv.Atoi(params["taskId"])
	if err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	tasks, createErr := service.CompleteTask(userId, projectId, taskId, h.todoRepository)
	if createErr != nil {
		h.SendError(w, createErr.StatusCode, createErr.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, tasks)
}
