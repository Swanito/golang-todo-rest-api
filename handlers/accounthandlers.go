package handlers

import (
	"encoding/json"
	"go-playground/model"
	"go-playground/service"
	"net/http"
)

func (h *BaseHandler) LoginHandler(w http.ResponseWriter, req *http.Request) {
	loginData := model.Login{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&loginData); err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := service.LoginService(loginData.Email, loginData.Password, h.userRepository)
	if err != nil {
		h.SendError(w, err.StatusCode, err.Error.Error())
		return
	}
	h.Send(w, http.StatusOK, user)

}

func (h *BaseHandler) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	newUser := model.NewUser{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&newUser); err != nil {
		h.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	errEmail := service.CheckEmail(newUser.Email, h.userRepository)
	if errEmail != nil {
		h.SendError(w, errEmail.StatusCode, errEmail.Error.Error())
		return
	}

	errUsername := service.CheckUsername(newUser.Username, h.userRepository)
	if errUsername != nil {
		h.SendError(w, errUsername.StatusCode, errUsername.Error.Error())
		return
	}

	userId, err := service.RegisterService(newUser.Username, newUser.Password, newUser.ConfirmPassword, newUser.Email, h.userRepository)
	if err != nil {
		h.SendError(w, err.StatusCode, err.Error.Error())
		return
	}

	emailErr := service.SendActivationEmail(newUser.Username, newUser.Email)
	if emailErr != nil {
		h.SendError(w, emailErr.StatusCode, emailErr.Error.Error())
		return
	}
	h.Send(w, http.StatusCreated, model.RegisterResponse{UserId: userId})
}

func (h *BaseHandler) ActivateAccountHandler(w http.ResponseWriter, req *http.Request) {
	userEmail := req.URL.Query().Get("email")
	if userEmail == "" {
		h.SendError(w, http.StatusNotFound, "No email found")
		return
	}
	err := service.ActivateUser(userEmail, h.userRepository)
	if err != nil {
		h.SendError(w, err.StatusCode, err.Error.Error())
		return
	}
	h.Send(w, http.StatusCreated, map[string]interface{}{"status": "Account activated successfully"})
}
