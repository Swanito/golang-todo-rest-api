package repository

import (
	"database/sql"
	"errors"
	"go-playground/model"

	"time"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (ar *AccountRepository) GetUserByUsername(name string) (*model.User, error) {
	var username, email, password string
	var lastLogin, activatedAt sql.NullTime
	var createdAt time.Time
	var id, loginAttempts int
	var activated bool
	err := ar.db.QueryRow(`SELECT * FROM USERS WHERE USERNAME = $1;`, name).Scan(&username, &email, &password, &lastLogin, &createdAt, &id, &loginAttempts, &activated, &activatedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(err.Error())
	}
	user := model.User{
		Id:            id,
		Email:         email,
		Username:      username,
		Password:      password,
		CreatedAt:     createdAt,
		LoginAttempts: loginAttempts,
	}
	return &user, nil
}

func (ar *AccountRepository) GetUserByEmail(userEmail string) (*model.User, error) {
	var username, email, password string
	var lastLogin, activatedAt sql.NullTime
	var createdAt time.Time
	var id, loginAttempts int
	var activated bool
	err := ar.db.QueryRow(`SELECT * FROM USERS WHERE EMAIL = $1;`, userEmail).Scan(&username, &email, &password, &lastLogin, &createdAt, &id, &loginAttempts, &activated, &activatedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(err.Error())
	}
	user := model.User{
		Id:            id,
		Email:         email,
		Username:      username,
		Password:      password,
		CreatedAt:     createdAt,
		LoginAttempts: loginAttempts,
		Activated:     activated,
	}
	return &user, nil
}

func (ar *AccountRepository) Register(username, password, email string, createdAt time.Time) (int, error) {
	var id int
	err := ar.db.QueryRow(`INSERT INTO USERS (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id;`, username, email, password, createdAt).Scan(&id)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return id, nil
}

func (ar *AccountRepository) GetLoginAttempts(email string) (int, error) {
	var loginAttempts int

	errSelect := ar.db.QueryRow(`SELECT LOGIN_ATTEMPTS FROM USERS WHERE EMAIL = $1`, email).Scan(&loginAttempts)
	if errSelect != nil && errSelect != sql.ErrNoRows {
		return 0, errors.New(errSelect.Error())
	}

	return loginAttempts, nil
}

func (ar *AccountRepository) AddInvalidLoginAttempt(email string) error {
	errUpdate := ar.db.QueryRow(`UPDATE USERS SET LOGIN_ATTEMPTS = LOGIN_ATTEMPTS + 1 WHERE EMAIL = $1`, email).Err()
	if errUpdate != nil && errUpdate != sql.ErrNoRows {
		return errors.New(errUpdate.Error())
	}

	return nil
}

func (ar *AccountRepository) AddLoginSession(id int, sessionDate time.Time) error {
	errInsert := ar.db.QueryRow(`INSERT INTO SESSIONS (user_id, session_date) VALUES ($1, $2)`, id, sessionDate).Err()
	if errInsert != nil {
		return errors.New(errInsert.Error())
	}
	return nil
}

func (ar *AccountRepository) ActivateUser(userEmail string) error {
	err := ar.db.QueryRow(`UPDATE USERS SET ACTIVATED = true, ACTIVATED_AT = $1 WHERE EMAIL = $2`, time.Now(), userEmail).Err()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
