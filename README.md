<h1 align="center">Тестовое Задание</h1>

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/8.png?raw=true)

<h1 align="center">Развертка</h1>

- Склонировать репозиторий
```
git clone https://github.com/Yury132/Golang-Task-4.git
```
- Установить PostgreSQL в Docker контейнер, используя docker-compose.yml файл из проекта
  
1. Скопировать docker-compose.yml в новую папку "postgresql"
  
2. Выполнить в терминале команду
```
docker compose up
```
- Подключиться к базе данных PostgreSQL (Например, через DBeaver)

POSTGRES_DB: mydb

POSTGRES_USER: root

POSTGRES_PASSWORD: mydbpass

Port: 5432

Host: localhost

- Запустить веб-приложение командой
```
go run cmd/main.go
```

<h1 align="center">Тестирование</h1>

- Используя браузер, перейти по следующему адресу

```
http://localhost:8080/users-list
```
Добавить нового пользователя, указав его ФИО

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/1.png?raw=true)

Нажать на кнопку "Добавить", пользователь отобразится в списке ниже

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/2.png?raw=true)

При нажатии на крестик запись о данном пользователе будет удалена из БД

При нажатии на ФИО пользователя отобразится следующий экран

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/3.png?raw=true)

Нажав на кнопку "Изменить ФИО пользователя", можно отредактировать его данные

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/4.png?raw=true)

Список добавленных пользователей можно фильтровать по возрасту, нажав на кнопку "Фильтр по возрасту"

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/5.png?raw=true)

Также доступен фильтр по полу - отображение только мужчин или женщин

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/6.png?raw=true)

Присутствует фильтр по национальности, где требуется указать код страны, например, "RU", "UA" и т.д.

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/7.png?raw=true)

При нажатии на кнопку "Сбросить фильтр" отображаются все пользователи системы без какой-либо дополнительной фильтрации


