package main

import (
	"fmt"
	"log"
)

type (
	AppName string
	DBName  string
)

type Logger struct {
	prefix AppName
}

func NewLogger(prefix AppName) *Logger {
	return &Logger{prefix}
}

func (l *Logger) Log(message string) {
	log.Printf("[%s] %s\n", l.prefix, message)
}

type UserRepository struct {
	dbName DBName
	logger *Logger
}

func NewUserRepository(dbName DBName, logger *Logger) *UserRepository {
	return &UserRepository{
		dbName: dbName,
		logger: logger,
	}
}

func (r *UserRepository) SaveUser(user string) {
	r.logger.Log(fmt.Sprintf("Saving user '%s' to database '%s'", user, r.dbName))
}

type UserService struct {
	repo   *UserRepository
	logger *Logger
}

func NewUserService(repo *UserRepository, logger *Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) RegisterUser(username string) {
	s.logger.Log(fmt.Sprintf("Attempting to register user: %s", username))
	s.repo.SaveUser(username)
}

type TaskRepository struct {
	dbName DBName
	logger *Logger
}

func NewTaskRepository(dbName DBName, logger *Logger) *TaskRepository {
	return &TaskRepository{
		dbName: dbName,
		logger: logger,
	}
}

func (r *TaskRepository) SaveTask(task string) {
	r.logger.Log(fmt.Sprintf("create task '%s' to database '%s'", task, r.dbName))
}

type TaskService struct {
	repo   *TaskRepository
	logger *Logger
}

func NewTaskService(repo *TaskRepository, logger *Logger) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}

func (t *TaskService) CreateTask(task string) {
	t.logger.Log(fmt.Sprintf("Attempting to create task: %s", task))
	t.repo.SaveTask(task)
}

func main() {
	userService := InitializeUserService()

	userService.RegisterUser("Alice")
}
