package main

import (
	"go-playground/handlers"
	"go-playground/repository"
	"go-playground/server"
)

func main() {
	app := server.GetApp()
	db, err := server.ConnectDatabase(app.DBConfig)
	if err != nil {
		panic(err)
	}
	accountRepository := repository.NewAccountRepository(db)
	todoRepository := repository.NewTodoRepository(db)
	handlers := handlers.NewBaseHandler(accountRepository, todoRepository)
	app.SetRoutes(handlers)
	app.Run(app.HTTPPort, app.Router)
}
