<h1 align="center">Тестовое Задание</h1>
![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/8.PNG?raw=true)
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

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/1.PNG?raw=true)

Нажать на кнопку "Добавить", пользователь отобразится в списке ниже

![alt text](https://github.com/Yury132/Golang-Task-4/blob/main/forREADME/2.PNG?raw=true)

