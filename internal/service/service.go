package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Yury132/Golang-Task-4/internal/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Service interface {
	// Получение всех пользователей БД
	GetUsersList(ctx context.Context) ([]models.User, error)
	// Получение определенных пользователей по возрасту
	GetUsersListAge(ctx context.Context, userAgeMin int, userAgeMax int) ([]models.User, error)
	// Получение определенных пользователей по полу
	GetUsersListGender(ctx context.Context, gender string) ([]models.User, error)
	// Получение определенных пользователей по национальности
	GetUsersListNation(ctx context.Context, userNation string) ([]models.User, error)
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
	// Получаем данные о возрасте
	GetAge(name string) ([]byte, error)
	// Получаем данные о поле
	GetGender(name string) ([]byte, error)
	// Получаем данные о национальности
	GetNation(name string) ([]byte, error)
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
	CreateUser(ctx context.Context, name string, surname string, patronymic string, age int, gender string, nation string) error
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

// Получение определенных пользователей по возрасту
func (s *service) GetUsersListAge(ctx context.Context, userAgeMin int, userAgeMax int) ([]models.User, error) {
	users, err := s.storage.GetUsersListAge(ctx, userAgeMin, userAgeMax)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Получение определенных пользователей по полу
func (s *service) GetUsersListGender(ctx context.Context, gender string) ([]models.User, error) {
	users, err := s.storage.GetUsersListGender(ctx, gender)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Получение определенных пользователей по национальности
func (s *service) GetUsersListNation(ctx context.Context, userNation string) ([]models.User, error) {
	users, err := s.storage.GetUsersListNation(ctx, userNation)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Добавление нового пользователя, если точно такой же уже не существует в БД
func (s *service) HandleUser(ctx context.Context, name string, surname string, patronymic string) error {

	// Используем api для получения возраста
	dataBytesAge, err := s.userAPI.GetAge(name)
	if err != nil {
		return errors.Wrap(err, "failed to get age from api")
	}

	// Переводим байты в структуру
	var infoAge models.AgeApi
	if err = json.Unmarshal(dataBytesAge, &infoAge); err != nil {
		return errors.Wrap(err, "failed to unmarshal age from api")
	}
	s.logger.Log().Msg(fmt.Sprintf("Для %v api вернул возраст %v", name, infoAge.Age))

	// Используем api для получения пола
	dataBytesGender, err := s.userAPI.GetGender(name)
	if err != nil {
		return errors.Wrap(err, "failed to get gender from api")
	}

	// Переводим байты в структуру
	var infoGender models.GenderApi
	if err = json.Unmarshal(dataBytesGender, &infoGender); err != nil {
		return errors.Wrap(err, "failed to unmarshal gender from api")
	}
	s.logger.Log().Msg(fmt.Sprintf("Для %v api вернул пол %v", name, infoGender.Gender))

	// Для БД формируем обозначение пол пользователя
	var getGender string
	if infoGender.Gender == "male" {
		getGender = "м"
	} else {
		getGender = "ж"
	}

	// Используем api для получения национальности
	dataBytesNation, err := s.userAPI.GetNation(name)
	if err != nil {
		return errors.Wrap(err, "failed to get nation from api")
	}

	// Переводим байты в структуру
	var infoNation models.NationApi
	if err = json.Unmarshal(dataBytesNation, &infoNation); err != nil {
		return errors.Wrap(err, "failed to unmarshal nation from api")
	}
	s.logger.Log().Msg(fmt.Sprintf("Для %v api вернул следующие коды стран: %v", name, infoNation.Country))

	// Ищем наибольшую вероятность - какую национальность имеет пользователь
	probability := 0.0
	countryCode := "RU"
	// Проходимся в цикле
	for _, value := range infoNation.Country {
		if value.Probability > probability {
			probability = value.Probability
			countryCode = value.Country_id
		}
	}

	// Проверяем на полное совпадение по ФИО в БД
	ok, err := s.checkUser(ctx, name, surname, patronymic)
	if err != nil {
		return errors.Wrap(err, "failed to check user")
	}

	if ok {
		s.logger.Log().Msg("ФИО нового пользователя полностью совпадает с уже существующим")
	} else {
		// Создаем
		if err = s.createUser(ctx, name, surname, patronymic, infoAge.Age, getGender, countryCode); err != nil { // infoAge.Age, getGender
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
func (s *service) createUser(ctx context.Context, name string, surname string, patronymic string, age int, gender string, nation string) error {
	err := s.storage.CreateUser(ctx, name, surname, patronymic, age, gender, nation)
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
