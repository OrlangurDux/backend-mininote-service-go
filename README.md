# Backend Service Mini Note
> backend microservice mini note golang version
## Установка
#### 1. Клонируем репозиторий
```bash
$ git clone https://github.com/OrlangurDux/backend-mininote-service-go.git
```

## Упаковка и Развертывание
#### Сборка и запуск без Docker
```bash
$ cd src
$ go build 
```

#### Сборка и запуск из Docker файла
```bash
$ docker run -d -v $(pwd)/mongodb:/data/db --name service-mini-note-mongo --network mini_note_network mongo:latest
$ docker build -t service-mini-note .
$ docker run -d -p 9077:9077 --name=service-mini-note --network mini_note_network -e "MONGODB_URL=mongodb://service-mini-note-mongo:27017/db" service-mini-note
```

#### Сборка и запуск с Docker-Compose
```
$ docker-compose up -d --build
```

## Разработка

### Запуск тестового сервера
Запускаем базу данных MongoDB как сервис в docker контейнере используя docker-compose скрипт `docker-compose.dev.yml`.

#### 1. Устанавливаем вендорные зависимости
```bash
$ go mod vendor
```

#### 2. Запуск тестов
```bash
$ cd tests
$ go test
```

#### 3. Запускаем приложение
```bash
$ go run main.go routes.go
```
Запускаем сервисы, результаты в
* 🌏**API Сервер** запущен на `http://localhost:9077`
* ⚙️**Swagger UI** запущен на `http://localhost:9077/swagger/index.html`

---

## Переменные окружения
Для редактирования переменных окружения, создайте файл с именем `.env` и скопируйте содержимое из `.env.default` чтобы начать.  
При запуске через Docker файл дополнительно можно указать перменные окружения используя агрумент `-e`  
При запуске через Docker-Compose можно указать перменные окружения через директиву `environment`

| Var Name           | Type     | Default                        | Description  |
|--------------------|----------|--------------------------------|---|
| JWT_SECRET         | string   | `secret`                       |JWT секрет для верификации  |
| JWT_LIFETIME       | number   | `720`                          | Время жизни JWT токена     |
| PORT               | number   | `9077`                         | Порт на котором запускается API сервер |
| UPDATE_IMPORT_HOUR | number   | `8`                            | Интервал обновления в часах |
| MONGODB_URL        | string   | `mongodb://localhost:27017/db` | Хост подключения к БД |
| AVATAR_DIR         | string   | `/uploaded/avatars/`           | Путь для сохранения аваторок профиля пользователя |
| HOST               | string   | `http://localhost:9077`        | Доменное имя проекта |
| DEFAULT_PER_PAGE   | number   | `20`                           | Количество элементов на страницу в списке |
| SMTP_HOST          | string   | `localhost`                    | Хост SMTP сервера |
| SMTP_LOGIN         | string   | `login`                        | Логин SMTP |
| SMTP_PASSWORD      | string   | `password`                     | Пароль SMTP |
| SMTP_PORT          | number   | `587`                          | Порт SMTP |
| SMTP_MESSAGE       | string   | `Message text #name# #phone#`  | Текст сообщения |
| SMTP_SUBJECT       | string   | `Subject text`                 | Тема сообщения |
| SMTP_TO            | string   | `to@mail.loc`                  | Получатель системных сообщений |