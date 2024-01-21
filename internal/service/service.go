package service

import (
	"context"

	"github.com/Yury132/Golang-Task-4/internal/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Service interface {
	// Получение всех пользователей БД
	GetUsersList(ctx context.Context) ([]models.User, error)
	// Добавление нового пользователя, если точно такой же уже не существует в БД
	HandleUser(ctx context.Context, name string, surname string, patronymic string) error
	// Удаление пользователя
	DeleteUser(ctx context.Context, id int) error
	// Получение конкретного пользователя по ID
	GetUser(ctx context.Context, id int) (models.User, error)
	// Обновление данных конкретного пользователя по ID
	EditUser(ctx context.Context, id int, getUserName string, getUserSurname string, getUserPatronymic string) error
}

type UserAPI interface {
	GetAge(name string) ([]byte, error)
}

type Storage interface {
	// Получение всех пользователей БД
	GetUsersList(ctx context.Context) ([]models.User, error)
	// Получение определенных пользователей по возрасту
	GetUsersListAge(ctx context.Context, ageMin int, ageMax int) ([]models.User, error)
	// Получение определенных пользователей по полу
	GetUsersListGender(ctx context.Context, gender string) ([]models.User, error)
	// Получение определенных пользователей по национальности
	GetUsersListNation(ctx context.Context, nation string) ([]models.User, error)
	// Проверка на существование пользователя
	CheckUser(ctx context.Context, name string, surname string, patronymic string) (bool, error)
	// Создание нового пользователя
	CreateUser(ctx context.Context, name string, surname string, patronymic string) error
	// Удаление пользователя
	DeleteUser(ctx context.Context, id int) error
	// Получение конкретного пользователя по ID
	GetUser(ctx context.Context, id int) (models.User, error)
	// Обновление данных конкретного пользователя по ID
	EditUser(ctx context.Context, id int, getUserName string, getUserSurname string, getUserPatronymic string) error
}

type service struct {
	logger  zerolog.Logger
	userAPI UserAPI
	storage Storage
}

// Все пользователи в БД
func (s *service) GetUsersList(ctx context.Context) ([]models.User, error) {
	users, err := s.storage.GetUsersList(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Добавление нового пользователя, если точно такой же уже не существует в БД
func (s *service) HandleUser(ctx context.Context, name string, surname string, patronymic string) error {
	// Проверяем
	ok, err := s.checkUser(ctx, name, surname, patronymic)
	if err != nil {
		return errors.Wrap(err, "failed to check user")
	}

	if ok {
		s.logger.Log().Msg("ФИО нового пользователя полностью совпадает с уже существующим")
	} else {
		// Создаем
		if err = s.createUser(ctx, name, surname, patronymic); err != nil {
			return errors.Wrap(err, "failed to create user")
		}
	}

	return nil
}

// Проверка на существование пользователя
func (s *service) checkUser(ctx context.Context, name string, surname string, patronymic string) (bool, error) {
	check, err := s.storage.CheckUser(ctx, name, surname, patronymic)
	if err != nil {
		return false, err
	}

	return check, nil
}

// Создание нового пользователя
func (s *service) createUser(ctx context.Context, name string, surname string, patronymic string) error {
	err := s.storage.CreateUser(ctx, name, surname, patronymic)
	if err != nil {
		return err
	}

	return nil
}

// Удаление пользователя
func (s *service) DeleteUser(ctx context.Context, id int) error {

	// Удаляем пользователя из БД
	err := s.storage.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Получение конкретного пользователя по ID
func (s *service) GetUser(ctx context.Context, id int) (models.User, error) {
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// Обновление данных конкретного пользователя по ID
func (s *service) EditUser(ctx context.Context, id int, getUserName string, getUserSurname string, getUserPatronymic string) error {
	err := s.storage.EditUser(ctx, id, getUserName, getUserSurname, getUserPatronymic)
	if err != nil {
		return err
	}

	return nil
}

func New(logger zerolog.Logger, userAPI UserAPI, storage Storage) Service {
	return &service{
		logger:  logger,
		userAPI: userAPI,
		storage: storage,
	}
}
