package server

import (
	"encoding/json"
	"go-playground/handlers"
	"go-playground/routes"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	HTTPPort string
	DBConfig DBConfig
	Router   *mux.Router
}

type DBConfig struct {
	Username string
	Password string
	DBName   string
	Host     string
	Port     string
}

func (app *App) SetRoutes(handlers *handlers.BaseHandler) *mux.Router {
	app.Router = mux.NewRouter()
	app.Router.HandleFunc(routes.LoginPath, app.interceptor(handlers.LoginHandler, false)).Name("LoginRoute").Methods("POST")
	app.Router.HandleFunc(routes.RegisterPath, app.interceptor(handlers.RegisterHandler, false)).Name("RegisterRoute").Methods("POST")
	app.Router.HandleFunc(routes.ActivateAccountPath, app.interceptor(handlers.ActivateAccountHandler, false)).Name("ActivateAccountRoute").Methods("GET")
	app.Router.HandleFunc(routes.ProjectsPath, app.interceptor(handlers.GetProjectsHandler, true)).Name("GetProjectsRoute").Methods("GET")
	app.Router.HandleFunc(routes.ProjectsPath, app.interceptor(handlers.CreateProjectHandler, true)).Name("CreateProjectsRoute").Methods("POST")
	app.Router.HandleFunc(routes.ProjectPath, app.interceptor(handlers.GetProjectHandler, true)).Name("GetProjectRoute").Methods("GET")
	app.Router.HandleFunc(routes.ProjectPath, app.interceptor(handlers.UpdateProjectHandler, true)).Name("UpdateProjectRoute").Methods("PUT")
	app.Router.HandleFunc(routes.ProjectPath, app.interceptor(handlers.DeleteProjectHandler, true)).Name("DeleteProjectRoute").Methods("DELETE")
	app.Router.HandleFunc(routes.ProjectPath, app.interceptor(handlers.CompleteProjectHandler, true)).Name("CompleteProjectRoute").Methods("POST")
	app.Router.HandleFunc(routes.TasksPath, app.interceptor(handlers.CreateTaskHandler, true)).Name("CreateTaskRoute").Methods("POST")
	app.Router.HandleFunc(routes.TaskPath, app.interceptor(handlers.CompleteTaskHandler, true)).Name("CompleteTaskRoute").Methods("POST")
	app.Router.HandleFunc(routes.TasksPath, app.interceptor(handlers.GetTasksHandler, true)).Name("GetTasksRoute").Methods("GET")
	return app.Router
}

func (app *App) Run(host string, router *mux.Router) {
	log.Fatal(http.ListenAndServe(":"+host, router))
}

func (app *App) interceptor(handler http.HandlerFunc, isPrivate bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if isPrivate {
			if req.Header["Authorization"] == nil {
				response, err := json.Marshal(map[string]interface{}{"error": "Unauthorized"})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(response))
				return
			}
			claims := handlers.ValidateToken(req.Header)
			userId := claims["userId"]
			req.Header.Add("userId", strconv.FormatInt(int64(userId.(float64)), 10))
		}
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, req)
	}
}
