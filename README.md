
# ToDo App (Golang CLI)

Консольное приложение для управления списком задач (todo) на языке Go.
Поддерживает команды add, list, delete, complete, import, export, гибкую фильтрацию и экспорт/импорт задач в форматах JSON/CSV.

---

## Возможности

- Добавление, выполнение и удаление задач
- Фильтрация задач: all, done, pending
- Экспорт и импорт задач в форматах JSON и CSV
- Умный CLI: используется командный стиль (add, list, complete и т.д.)
- Хранение задач в tasks.json (с автосохранением)

---

## Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Evgenymoshrage/todo-app.git
   cd todo-app
   ```

2. Собрать исполняемый файл:
   ```bash
   go build -o todo ./cmd/todo
   ```

3. Запуск:
   ```bash

   ./todo

   ```

## Использование

### Через флаги:

```bash
./todo add --desc="Прочитать книгу"
./todo list --filter=pending
./todo complete --id=2
./todo delete --id=3
./todo export --file=backup.json --format=json
./todo import --file=backup.csv --format=csv
```

### Справка:

```bash

./todo --help

./todo -h
```

## Структура проекта

```bash
todo-app/
│
├── cmd/todo/                                # Точка входа
│        ├── main.go                         # Главный исходный файл Go-программы
│        ├── main_test.go                    # Содержит набор модульных тестов
│        ├── tasks.json                      # Файл с сохранёнными задачами
│        └── todo.exe                        # Скомпилированный исполняемый файл программы для Windows
├── internal/
│   ├── todo/                                # Управление задачами в памяти
│        ├── manager/manager.go              # Главный исходный файл Go-программы
│        └── task/task.go                    # Файл с сохранёнными задачами
│   └── storage/
│       ├── csv_storage/csv_storage.go       # Импорт/экспорт CSV
│       └── json_storage/json_storage.go     # Импорт/экспорт JSON
├── go.mod                                   # Файл управления зависимостями в Go
└── README.md                                # Основной файл документации проекта
```

## Тестирование приложения (main_test.go)
```bash
Файл main_test.go содержит набор модульных тестов, предназначенных для проверки корректной работы основных функций todo-приложения.

Что тестируется:
- **TestAddTask** — проверка добавления задачи.
- **TestDeleteTask** — проверка удаления задачи по ID.
- **TestCompleteTask** — проверка отметки задачи как выполненной.
- **TestListWithFilters** — проверка фильтрации задач (все, выполненные, невыполненные).
- **TestLoadJSON_FileNotFound** — проверка ошибки при загрузке несуществующего JSON-файла.
- **TestSaveAndLoadJSON_Success** — проверка корректного сохранения и загрузки задач в JSON.
- **TestLoadCSV_FileNotFound** — проверка ошибки при загрузке несуществующего CSV-файла.
- **TestSaveAndLoadCSV_Success** — проверка корректного сохранения и загрузки задач в CSV.
```
### Особенности:

Каждый тест очищает список задач перед запуском с помощью manager.ClearTasks(), чтобы избежать конфликта с предыдущими данными.

Проверки выполняются с использованием стандартной библиотеки testing.

## Как запустить тест:

```bash

go test ./... -v

```

[Evgenymoshrage](https://github.com/Evgenymoshrage)



