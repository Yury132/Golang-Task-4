package storage

import (
	"context"

	"github.com/Yury132/Golang-Task-4/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

type storage struct {
	conn *pgxpool.Pool
}

// Все пользователи в БД
func (s *storage) GetUsersList(ctx context.Context) ([]models.User, error) {
	query := "SELECT id, name, surname, patronymic, age, gender, nation FROM public.users"

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nation); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

// Получение определенных пользователей по возрасту
func (s *storage) GetUsersListAge(ctx context.Context, ageMin int, ageMax int) ([]models.User, error) {
	query := "SELECT id, name, surname, patronymic, age, gender, nation FROM public.users WHERE age >= $1 AND age <= $2"

	rows, err := s.conn.Query(ctx, query, ageMin, ageMax)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nation); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

// Получение определенных пользователей по полу
func (s *storage) GetUsersListGender(ctx context.Context, gender string) ([]models.User, error) {
	query := "SELECT id, name, surname, patronymic, age, gender, nation FROM public.users WHERE gender = $1"

	rows, err := s.conn.Query(ctx, query, gender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nation); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

// Получение определенных пользователей по национальности
func (s *storage) GetUsersListNation(ctx context.Context, nation string) ([]models.User, error) {
	query := "SELECT id, name, surname, patronymic, age, gender, nation FROM public.users WHERE nation = $1"

	rows, err := s.conn.Query(ctx, query, nation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nation); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

// Проверка на существование пользователя
func (s *storage) CheckUser(ctx context.Context, name string, surname string, patronymic string) (bool, error) {
	query := "SELECT id FROM public.users WHERE name = $1 AND surname = $2 AND patronymic = $3"

	rows, err := s.conn.Query(ctx, query, name, surname, patronymic)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Считаем количество строк
	countRows := 0
	for rows.Next() {
		countRows++
	}
	// Если что то нашли, значит пользователь есть в БД
	check := false
	if countRows > 0 {
		check = true
	}

	if rows.Err() != nil {
		return false, err
	}

	return check, nil
}

// Создание нового пользователя
func (s *storage) CreateUser(ctx context.Context, name string, surname string, patronymic string, age int, gender string, nation string) error {
	query := "INSERT INTO public.users (name, surname, patronymic, age, gender, nation) values ($1, $2, $3, $4, $5, $6)"
	_, err := s.conn.Exec(ctx, query, name, surname, patronymic, age, gender, nation)
	if err != nil {
		return err
	}
	return nil
}

// Удаление пользователя
func (s *storage) DeleteUser(ctx context.Context, id int) error {
	query := "delete from public.users where id = $1"

	_, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// Получение конкретного пользователя по ID
func (s *storage) GetUser(ctx context.Context, id int) (models.User, error) {
	// Структура
	var user models.User
	// Запрос - Получаем только одну строку
	query := "SELECT id, name, surname, patronymic, age, gender, nation FROM public.users WHERE id = $1"
	// Выполняем запрос, возвращающий только одну строку
	row := s.conn.QueryRow(ctx, query, id)

	// Считываем значение
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nation); err != nil {
		return user, err
	}
	return user, nil
}

// Обновление данных конкретного пользователя по ID
func (s *storage) EditUser(ctx context.Context, id int, getUserName string, getUserSurname string, getUserPatronymic string) error {
	query := "UPDATE public.users SET name = $1, surname = $2, patronymic = $3 WHERE id = $4"
	_, err := s.conn.Exec(ctx, query, getUserName, getUserSurname, getUserPatronymic, id)
	if err != nil {
		return err
	}
	return nil
}

func New(conn *pgxpool.Pool) Storage {
	return &storage{
		conn: conn,
	}
}
