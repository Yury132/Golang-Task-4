package service

import (
	"context"

	"github.com/Yury132/Golang-Task-4/internal/client/api"
	"github.com/Yury132/Golang-Task-4/internal/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Service interface {
	GetUsersList(ctx context.Context) ([]models.User, error)
	HandleUser(ctx context.Context, name string, email string) error
}

type UserAPI interface {
	GetAge(name string) ([]byte, error)
}

type Storage interface {
	// Получение
	GetUsers(ctx context.Context) ([]models.User, error)
	// Проверка на существование пользователя
	CheckUser(ctx context.Context, email string) (bool, error)
	// Создание нового пользователя
	CreateUser(ctx context.Context, name string, email string) error
}

type service struct {
	logger  zerolog.Logger
	userAPI UserAPI
	storage Storage
}

// Все пользователи в БД
func (s *service) GetUsersList(ctx context.Context) ([]models.User, error) {
	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Проверка существования пользователя в БД и его создание при необходимости
func (s *service) HandleUser(ctx context.Context, name string, email string) error {
	ok, err := s.checkUser(ctx, email)
	if err != nil {
		return errors.Wrap(err, "failed to check user")
	}

	if !ok {
		if err = s.createUser(ctx, name, email); err != nil {
			return errors.Wrap(err, "failed to create user")
		}
	}

	return nil
}

// Проверка на существование пользователя
func (s *service) checkUser(ctx context.Context, email string) (bool, error) {
	check, err := s.storage.CheckUser(ctx, email)
	if err != nil {
		return false, err
	}

	return check, nil
}

// Создание нового пользователя
func (s *service) createUser(ctx context.Context, name string, email string) error {
	err := s.storage.CreateUser(ctx, name, email)
	if err != nil {
		return err
	}

	return nil
}

func New(logger zerolog.Logger, userAPI api.UserAPI, storage Storage) Service {
	return &service{
		logger:  logger,
		userAPI: userAPI,
		storage: storage,
	}
}
