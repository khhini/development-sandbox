//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

var appProviderSet = wire.NewSet(
	wire.Value(AppName("APP")),
	wire.Value(DBName("user_db")),
	NewLogger,
	NewUserRepository,
	NewUserService,
	NewTaskRepository,
	NewTaskService,
)

func InitializeUserService() *UserService {
	wire.Build(appProviderSet)
	return &UserService{}
}

func InitializeTaskService() *TaskService {
	wire.Build(appProviderSet)
	return &TaskService{}
}
