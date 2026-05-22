Описание проекта

RESTful API для управления каталогом тротуарной плитки с полной системой аутентификации и авторизации. Реализовано в рамках лабораторной работы №3.

### Основные возможности

-  Регистрация и аутентификация пользователей
-  JWT токены (Access + Refresh) с ограниченным временем жизни
-  HttpOnly cookies для безопасной передачи токенов
-  Хеширование паролей с уникальной солью (bcrypt)
-  OAuth 2.0 (Yandex ID) с ручной реализацией Authorization Code Flow
-  Управление сессиями (Logout, Logout-all)
-  Защита CRUD операций из ЛР №2

Установка и запуск

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/Xeatric/paving-tiles-api.git
cd paving-tiles-api

cp .env.example .env

Отредактируйте .env

Запустите приложение:

docker-compose up --build

Регистрация пользователя
Цель: Создать нового пользователя в системе

Параметр	Значение
Method	POST
URL	{{baseUrl}}/auth/register
Headers	Content-Type: application/json
Body (raw JSON):

json
{
    "email": "{{user_email}}",
    "password": "{{user_password}}",
    "name": "{{user_name}}"
}
Логин (получение токенов)
Цель: Аутентифицировать пользователя и получить токены в cookies

Параметр	Значение
Method	POST
URL	{{baseUrl}}/auth/login
Headers	Content-Type: application/json
Body (raw JSON):
{
    "email": "{{user_email}}",
    "password": "{{user_password}}"
}

Whoami (проверка статуса)
Цель: Проверить, что пользователь авторизован

Параметр	Значение
Method	GET
URL	{{baseUrl}}/api/v1/auth/whoami
Headers	Authorization: Bearer {{access_token}}


Whoami без токена (отрицательный тест)
Цель: Проверить, что без токена доступ запрещён

Параметр	Значение
Method	GET
URL	{{baseUrl}}/api/v1/auth/whoami
Headers	(без Authorization)

Получение списка плиток (с токеном)
Цель: Получить защищённый ресурс с авторизацией

Параметр	Значение
Method	GET
URL	{{baseUrl}}/api/v1/tiles?page=1&limit=10
Headers	Authorization: Bearer {{access_token}}


Получение списка плиток (без токена)
Цель: Проверить, что без токена доступ запрещён

Параметр	Значение
Method	GET
URL	{{baseUrl}}/api/v1/tiles?page=1&limit=10
Headers	(без Authorization)