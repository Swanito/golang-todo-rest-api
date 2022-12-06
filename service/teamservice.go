package service

import (
	"errors"
	"go-playground/model"
	"net/http"
)

func GetTeams(userId int, repository model.TeamRepository) (*model.Teams, *model.CustomError) {
	teams, err := repository.GetTeams(userId)
	if err != nil {
		return nil, &model.CustomError{Error: errors.New(err.Error()), StatusCode: http.StatusInternalServerError}
	}
	return &model.Teams{Teams: teams}, nil
}
