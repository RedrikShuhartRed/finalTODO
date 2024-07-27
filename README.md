# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.
# Go веб-сервер планировщика задач

## Реализация
Реализованы все основные функции:
1. Запуск веб-сервиса
2. Создание БД
3. Правила повторения задач
4. Добавление задач
5. Получение списка задач
6. Редактирование задач
7. Отметка выполнения задач
8. Удаление задач
   
Реализованы все дополнительные функции (под *)
1. Определение переменных окружения из файла .env
2. Правила повторения под *
3. Поиск задач при передаче параметра search
4. Аутентификация
5. Создание Docker образа
   
## Локальный запуск
Из директории проекта ./cmd выполнить команду go run main.go:
```
$ go run cmd/main.go
```
Переменные для запуска по умолчанию(в файле .env):
```
TODO_PORT=7540
TODO_DBFILE=./scheduler.db
TODO_PASSWORD=myPassword
TODO_PASSWORDSALT = kl4509dafh43589whfh
TODO_TOKENSALT = klajglk54adgagsd
```
Перейти в браузере по адресу:
```
http://localhost:7540/login.html
```
Ввести пароль myPassword
## Тестирование
Для тестирования в основной директории проекта выполнить:
```
go clean -testcache
go test ./tests
```
Настроки переменных для тестирования в папке tests/settings, значения по умолчанию:
```
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoYXNoUGFzcyI6IjJhMzcxMzI1YjI1ZDQ1MDJlNzgyYzk2MTJiZmE0YTQ2MDE2ZjQxMzJhNGNjODllM2MyYWJkOTEwYjkxMzM5MGEifQ.l7vnj9evbxf_3GzYhcJ67Mvt-Ob4YKf7oVqHfH5Jl2o"
```
Для запуска отдельных тестов используйте:
```
go clean -testcache
go test -count=1 -run ^TestApp$ ./tests
go test -count=1 -run ^TestDB$ ./tests
go test -count=1 -run ^TestNextDate$ ./tests
go test -count=1 -run ^TestAddTask$ ./tests
go test -count=1 -run ^TestTasks$ ./tests
go test -count=1 -run ^TestTask$ ./tests
go test -count=1 -run ^TestEditTask$ ./tests
go test -count=1 -run ^TestDone$ ./tests
go test -count=1 -run ^TestDelTask$ ./tests
```
## Сборка Docker образа
В корневой директории проекта выполнить:
```$ docker build -t todo_server .```
## Запуск Docker контейнера
Для запуска контейнера в коневой директории проекта выполнить:
```docker run -d -p 8080:8080 -v /path/to/your/sqlite.db:/cmd/todo.db --name my-running-server todo_server```
Где /path/to/your/sqlite.db путь к вашей локальной SQLite, например:
```docker run -d -p 7540:7540 -v D:/sqllite:/cmd/todo.db --name my-running-server todo_server```
Для проверки работы контейнера перейти в браузера под адресу http://localhost:7540/
