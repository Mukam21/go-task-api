## Task Management REST API

Проект реализует REST API для управления задачами на Go с асинхронным логированием и in-memory хранилищем. Ниже представлены инструкции по запуску и тестированию API с помощью Postman.

Особенности реализации

1. Асинхронное логирование через канал и отдельную горутину

2. Потокобезопасное хранилище на sync.RWMutex и map

3. Graceful shutdown с обработкой SIGINT/SIGTERM

4. Валидация входных данных для создания задач

5. Фильтрация задач по статусу

# Запуск приложения

1. Убедитесь, что установлен Go (версия 1.20+)

2. Склонируйте репозиторий:

    git clone https://github.com/yourusername/go-task-api.git

    cd go-task-api

3. Запустите сервер:
    go run cmd/main.go

# Сервер запустится на порту 8080:
    2025/08/14 12:00:00 Server is running on :8080

2. Создание задачи (POST /tasks)

 - Выберите метод POST

 - Введите URL: http://localhost:8080/tasks

 - Перейдите в раздел Body → raw → JSON

 - Введите данные задачи:

    {
        "title": "Write documentation",
        "status": "in_progress"
    }

Ожидаемый результат:
    Статус 201 Created

    Тело ответа с созданной задачей:

    {
        "id": 1,
        "title": "Write documentation",
        "status": "in_progress",
        "created_at": "2025-08-14T12:05:00Z",
        "updated_at": "2025-08-14T12:05:00Z"
    }

3. Получение всех задач (GET /tasks)
 - Выберите метод GET

 - Введите URL: http://localhost:8080/tasks

 - Для фильтрации добавьте параметр:

    - Key: status

    - Value: in_progress

Ожидаемый результат:

 - Статус 200 OK

 - Список задач в формате JSON

4. Получение задачи по ID (GET /tasks/{id})

 - Выберите метод GET

 - Введите URL: http://localhost:8080/tasks/1

Ожидаемый результат:
 - Статус 200 OK
 - Данные задачи с ID=1

Ошибки:

 - 404 Not Found: Если задача не существует

 - 400 Bad Request: Если ID не число

Примеры запросов

# Успешное создание задачи

    POST /tasks HTTP/1.1
    Host: localhost:8080
    Content-Type: application/json

    {
        "title": "Review code",
        "status": "pending"
    }

Ответ

    HTTP/1.1 201 Created
    Content-Type: application/json

    {
        "id": 2,
        "title": "Review code",
        "status": "pending",
        "created_at": "2025-08-14T12:10:00Z",
        "updated_at": "2025-08-14T12:10:00Z"
    }

Фильтрация задач

    GET /tasks?status=pending HTTP/1.1
    Host: localhost:8080

Ответ

    HTTP/1.1 200 OK
    Content-Type: application/json

    [
        {
            "id": 2,
            "title": "Review code",
            "status": "pending",
            "created_at": "2025-08-14T12:10:00Z",
            "updated_at": "2025-08-14T12:10:00Z"
        }
    ]

Логи сервера

При выполнении запросов сервер выводит логи:

    2025-08-14T12:05:00Z [CreateTask] id=1 title="Write documentation" status=in_progress
    2025-08-14T12:10:00Z [CreateTask] id=2 title="Review code" status=pending
    2025-08-14T12:11:00Z [GetTasks] count=1 status=pending
    2025-08-14T12:12:00Z [GetTask] id=1

Graceful Shutdown

При завершении работы (Ctrl+C):

    2025/08/14 12:15:00 Shutting down...
    2025/08/14 12:15:00 Server stopped gracefully

# Технические детали реализации

 - Асинхронное логирование

    func (l *Logger) Start() {

        l.wg.Add(1)

        go func() {

            defer l.wg.Done()

            for msg := range l.ch {

                log.Printf("%s %s", time.Now().UTC().Format(time.RFC3339), msg)

            }

        }()

    }

 - Потокобезопасное хранилище

    func (r *TaskRepository) Create(task *model.Task) *model.Task {

        r.mu.Lock()

        defer r.mu.Unlock()

        r.lastID++

        task.ID = r.lastID

        task.CreatedAt = time.Now().UTC()

        task.UpdatedAt = task.CreatedAt

        r.tasks[task.ID] = task

        return task

    }

 - Graceful Shutdown

        stop := make(chan os.Signal, 1)

        signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

        <-stop

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

        defer cancel()

        srv.Shutdown(ctx)

        asyncLogger.Close()