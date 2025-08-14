## Требования

* Go **1.22+**

## Структура проекта

```
.
├── cmd/
│   └── main.go           # точка входа (HTTP-сервер)
└── internal/
    ├── handlers/         # HTTP-обработчики
    ├── logger/           # логгер
    ├── models/           # доменные структуры
    ├── services/         # бизнес-логика (use case)
    └── storage/          # in-memory репозиторий
```

---

## Запуск

```bash
git clone <repo-url> test-task-lo
cd test-task-lo

go run ./cmd
```

Остановка: `Ctrl+C`

---

## Сборка бинарников (Linux / macOS / Windows)

### Linux (amd64)

```bash
GOOS=linux GOARCH=amd64 go build -o bin/tasklo-linux-amd64 ./cmd
```

### macOS (amd64 / arm64)

```bash
GOOS=darwin GOARCH=amd64 go build -o bin/tasklo-darwin-amd64 ./cmd
GOOS=darwin GOARCH=arm64 go build -o bin/tasklo-darwin-arm64 ./cmd
```

### Windows (amd64)

```bash
GOOS=windows GOARCH=amd64 go build -o bin/tasklo-windows-amd64.exe ./cmd
```

---

## Запуск на разных ОС

### Linux / macOS (bash/zsh)

```bash
./bin/tasklo   
```

### Windows (PowerShell)

```powershell
.\bin\tasklo-windows-amd64.exe
```

Файрволл/антивирус может спросить разрешение на входящие подключения к порту `:8080`.

---

## API

Базовый URL: `http://localhost:8080`

### Схема

* `POST /tasks` – создать задачу
* `GET  /tasks/{id}` – получить задачу по ID
* `GET  /tasks?status=<pending|in_progress|completed>` – список задач (с фильтром по статусу опционально)

### Модель

```json
{
  "id": 1,
  "title": "string",
  "status": "pending | in_progress | completed",
  "created_at": "RFC3339 datetime",
  "updated_at": "RFC3339 datetime"
}
```

### Ответ об ошибке (JSON)

```json
{ "error": "message" }
```

Коды: `400` (валид. ошибка), `404` (не найдено), `500` (внутренняя).

---

## Примеры curl (Linux/macOS)

### Создать задачу

```bash
curl -s -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Buy milk","status":"pending"}' | jq .
```

### Получить по ID (пример: 1)

```bash
curl -s http://localhost:8080/tasks/1 
```

### Список всех задач

```bash
curl -s http://localhost:8080/tasks 
```

### Список задач со статусом `in_progress`

```bash
curl -s 'http://localhost:8080/tasks?status=in_progress'
```

---
