package model

type Teams struct {
	Teams []*Team
}

type Team struct {
	Name     string
	Id       int
	UserId   int
	Editable bool
}

type TeamRepository interface {
	GetTeams(userId int) ([]*Team, error)
}
