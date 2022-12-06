package server

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	app = App{
		Router: mux.NewRouter(),
	}
)

func init() {
	app.SetRoutes(nil)
}

func TestRoutes(t *testing.T) {
	t.Parallel()

	routesTable := []struct {
		name           string
		expectedMethod string
		expectedPath   string
	}{
		{"LoginRoute", http.MethodPost, "/login"},
		{"RegisterRoute", http.MethodPost, "/register"},
		{"ActivateAccountRoute", http.MethodGet, "/activate"},
		{"GetProjectsRoute", http.MethodGet, "/projects"},
		{"CreateProjectsRoute", http.MethodPost, "/projects"},
		{"GetProjectRoute", http.MethodGet, "/projects/{id}"},
		{"UpdateProjectRoute", http.MethodPut, "/projects/{id}"},
		{"DeleteProjectRoute", http.MethodDelete, "/projects/{id}"},
		{"CompleteProjectRoute", http.MethodPost, "/projects/{id}"},
		{"CreateTaskRoute", http.MethodPost, "/projects/{projectId}/tasks"},
		{"CompleteTaskRoute", http.MethodPost, "/projects/{projectId}/tasks/{taskId}"},
		{"GetTasksRoute", http.MethodGet, "/projects/{projectId}/tasks"},
	}
	for _, row := range routesTable {
		t.Run(row.name, func(t *testing.T) {
			// run
			route := app.Router.GetRoute(row.name)

			// verify
			if assert.NotNil(t, route) {
				methods, _ := route.GetMethods()
				assert.Equal(t, 1, len(methods))
				method := methods[0]
				assert.Equal(t, row.expectedMethod, method)
				path, _ := route.GetPathTemplate()
				assert.Equal(t, row.expectedPath, path)
			}
		})
	}
}

func TestInterceptorHeaders(t *testing.T) {
	t.Parallel()

	routesTable := []struct {
		isPrivate  bool
		authHeader string
	}{
		{false, ""},
		{true, ""},
		{true, "test"},
	}

	w := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "", nil)
	fn := func(_ http.ResponseWriter, _ *http.Request) {}

	for _, row := range routesTable {
		w.Header().Set("Authorization", row.authHeader)
		hndlFn := app.interceptor(fn, row.isPrivate)
		hndlFn.ServeHTTP(w, request)
		assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
		if row.isPrivate {
			assert.Equal(t, 401, w.Code)
		}
	}

}
