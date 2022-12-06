package repository

import (
	"database/sql"
	"errors"
	"go-playground/model"
	"time"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

func (tr *TodoRepository) GetProjects(userId int) ([]*model.ProjectResponse, error) {
	var t string
	var completed bool
	var completedAt sql.NullTime
	var projectId, user int
	var projects = make([]*model.ProjectResponse, 0)
	rows, err := tr.db.Query(`SELECT projects.* 
	FROM projects 
	JOIN users_projects
	ON projects.id = users_projects.project_id
	JOIN users
	ON users.id = users_projects.user_id
	WHERE users.id = $1`, userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for rows.Next() {
		rows.Scan(&user, &completed, &completedAt, &t, &projectId)
		projects = append(projects, &model.ProjectResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, UserId: user, Id: projectId})
	}
	return projects, nil
}

func (tr *TodoRepository) GetProject(userId int, projectId int) (*model.ProjectResponse, error) {
	var t string
	var completed bool
	var completedAt sql.NullTime
	var user int

	pErr := tr.db.QueryRow(`SELECT projects.* 
	FROM projects 
	JOIN users_projects
	ON projects.id = users_projects.project_id
	JOIN users
	ON users.id = users_projects.user_id
	WHERE users.id = $1
	AND projects.id = $2`, userId, projectId).Scan(&user, &completed, &completedAt, &t, &projectId)
	if pErr != nil && pErr != sql.ErrNoRows {
		return nil, errors.New(pErr.Error())
	}
	if pErr == sql.ErrNoRows {
		return &model.ProjectResponse{}, nil
	}
	tasks, tErr := tr.GetTasks(userId, projectId)
	if tErr != nil {
		return nil, errors.New(pErr.Error())
	}
	return &model.ProjectResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, Id: projectId, UserId: userId, Tasks: tasks}, nil
}

func (tr *TodoRepository) CreateProject(userId int, title string) (*model.ProjectResponse, error) {
	var t string
	var completed bool
	var completedAt sql.NullTime
	var projectId, user int
	err := tr.db.QueryRow(`INSERT INTO PROJECTS (USER_ID, TITLE, COMPLETED) VALUES ($1, $2, $3) RETURNING *`, userId, title, false).Scan(&user, &completed, &completedAt, &t, &projectId)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	teamErr := tr.db.QueryRow(`INSERT INTO USERS_PROJECTS (USER_ID, PROJECT_ID) VALUES ($1, $2)`, userId, projectId).Err()
	if teamErr != nil {
		return nil, errors.New(teamErr.Error())
	}
	return &model.ProjectResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, Id: projectId, UserId: userId}, nil
}

func (tr *TodoRepository) UpdateProject(userId int, projectId int, newTitle string) (*model.ProjectResponse, error) {
	var t string
	var completed bool
	var completedAt sql.NullTime
	err := tr.db.QueryRow(`UPDATE PROJECTS SET TITLE = $1 WHERE USER_ID = $2 AND ID = $3 RETURNING *`, newTitle, userId, projectId).Scan(&userId, &completed, &completedAt, &t, &projectId)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(err.Error())
	}
	if err == sql.ErrNoRows {
		return nil, errors.New("Could not update project")
	}
	return &model.ProjectResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, Id: projectId}, nil
}

func (tr *TodoRepository) CompleteProject(userId int, projectId int) (*model.ProjectResponse, error) {
	var t string
	var completed bool
	var completedAt sql.NullTime
	err := tr.db.QueryRow(`UPDATE PROJECTS SET COMPLETED = true, COMPLETED_AT = $1 WHERE USER_ID = $2 AND ID = $3 RETURNING *`, time.Now(), userId, projectId).Scan(&userId, &completed, &completedAt, &t, &projectId)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(err.Error())
	}
	if err == sql.ErrNoRows {
		return nil, errors.New("Could not update project")
	}
	return &model.ProjectResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, Id: projectId}, nil
}

func (tr *TodoRepository) DeleteProject(userId int, projectId int) error {
	errProject := tr.db.QueryRow(`DELETE FROM PROJECTS WHERE USER_ID = $1 AND ID = $2`, userId, projectId).Err()
	if errProject != nil {
		return errors.New(errProject.Error())
	}
	taskErr := tr.db.QueryRow(`DELETE FROM TASKS WHERE PROJECT_ID = $1`, projectId).Err()
	if taskErr != nil {
		return errors.New(taskErr.Error())
	}
	return nil
}

func (tr *TodoRepository) CreateTask(userId, projectId int, title, description string) (*model.TaskResponse, error) {
	var t, desc string
	var completed bool
	var completedAt sql.NullTime
	var id int
	err := tr.db.QueryRow(`INSERT INTO TASKS (COMPLETED, COMPLETED_AT, TITLE, DESCRIPTION, PROJECT_ID, USER_ID) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`, false, nil, title, description, projectId, userId).Scan(&id, &completed, &completedAt, &t, &desc, &projectId, &userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &model.TaskResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, Id: id, Description: desc, UserId: userId}, nil
}

func (tr *TodoRepository) GetTasks(userId, projectId int) ([]*model.TaskResponse, error) {
	var t, desc string
	var completed bool
	var completedAt sql.NullTime
	var id int
	var tasks = make([]*model.TaskResponse, 0)
	rows, err := tr.db.Query(`SELECT tasks.* 
	FROM tasks
	where tasks.project_id = (
		select users_projects.project_id 
		from users_projects 
		where users_projects.user_id = $1 and users_projects.project_id = $2)`, userId, projectId)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for rows.Next() {
		rows.Scan(&id, &completed, &completedAt, &t, &desc, &projectId, &userId)
		tasks = append(tasks, &model.TaskResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, Id: id, Description: desc, UserId: userId})
	}
	return tasks, nil
}

func (tr *TodoRepository) CompleteTask(userId, projectId, taskId int) (*model.TaskResponse, error) {
	var t, desc string
	var completed bool
	var completedAt sql.NullTime
	var id int
	err := tr.db.QueryRow(`UPDATE TASKS SET COMPLETED = true, COMPLETED_AT = $1 WHERE USER_ID = $2 AND PROJECT_ID = $3 AND ID = $4 RETURNING *`, time.Now(), userId, projectId, taskId).Scan(&id, &completed, &completedAt, &t, &desc, &projectId, &userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &model.TaskResponse{Title: t, Completed: completed, CompletedAt: completedAt.Time, Id: id, Description: desc, UserId: userId}, nil
}
