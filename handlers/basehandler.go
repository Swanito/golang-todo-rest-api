package handlers

import (
	"encoding/json"
	"go-playground/model"

	"net/http"
)

type BaseHandler struct {
	userRepository model.UserRepository
	todoRepository model.TodoRepository
}

func NewBaseHandler(userRepository model.UserRepository, todoRepository model.TodoRepository) *BaseHandler {
	return &BaseHandler{
		userRepository: userRepository,
		todoRepository: todoRepository,
	}
}

func (bh *BaseHandler) Send(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

func (bh *BaseHandler) SendError(w http.ResponseWriter, statusCode int, msg string) {
	bh.Send(w, statusCode, model.ApiError{Message: msg, StatusCode: statusCode})
}
