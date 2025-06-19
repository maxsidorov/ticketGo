## **ticketGo \- платформа для покупки билетов на мероприятия**

### Ссылка на сайт: not in work

### Команда:

1. Сидоров Максим  
2. Яшин Кирилл  
3. Дингилевский Михаил  
4. Моисеев Михаил 

### Технологический стек:

* git  
* Golang  
* Gin  
* PostgreSQL
* Docker
* Docker Compose

### Архитектура: Монолит

### Кол-во Docker контейнеров: 2

### Функционал:

1. Вывод доступных мероприятий, с возможностью поиска/фильтрации)  
2. Авторизации  
3. Регистрации  
4. Доп для админа:  
* Добавлением мероприятия  
* Возможность редактирования/удаления мероприятий  
5. Покупка билетов  
6. Генерация билетов (QR-код/PDF-файл)  
7. Возможность экспорта данных о мероприятиях в JSON

### Структура приложения:

* На каждой странице navbar  
* Главная \- страница с доступными мероприятиями  
* Страницы авторизации/регистрации  
* Страница с описанием мероприятия  
* Страница для покупки билета (с возможностью скачивания)  
* Для админа: интерфейс для добавления/удаления \+ редактирования мероприятий

---

## 🚀 Быстрый старт

### Локальный запуск

1. Установите PostgreSQL и создайте базу данных `postgres` с пользователем `postgres` и паролем `postgres`.

2. В файле `config.go` измените `DBPort` на тот, что указали при создании БД.

3. Запустите приложение:

   ```sh
   go run cmd/app/main.go
   ```

### Запуск с Docker

1. Скопируйте файл конфигурации:
   ```bash
   cp env.example .env
   ```

2. Отредактируйте `.env` файл с вашими настройками

3. Запустите приложение:
   ```bash
   docker-compose up -d
   ```

---

## 🌐 Развертывание на VPS/VDS

Подробное руководство по развертыванию на сервере находится в файле [DEPLOYMENT.md](DEPLOYMENT.md).

### Краткая инструкция:

1. **Подготовка сервера:**
   ```bash
   # Обновление системы
   sudo apt update && sudo apt upgrade -y
   
   # Установка Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   
   # Установка Docker Compose
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   ```

2. **Развертывание приложения:**
   ```bash
   # Клонирование репозитория
   sudo mkdir -p /opt/ticketgo
   sudo chown $USER:$USER /opt/ticketgo
   cd /opt/ticketgo
   git clone https://github.com/your-username/ticketGo.git .
   
   # Настройка переменных окружения
   cp env.example .env
   # Отредактируйте .env файл
   
   # Запуск приложения
   chmod +x deploy.sh
   ./deploy.sh
   ```

3. **Настройка Nginx (опционально):**
   ```bash
   sudo apt install nginx -y
   sudo cp nginx.conf /etc/nginx/sites-available/ticketgo
   sudo ln -s /etc/nginx/sites-available/ticketgo /etc/nginx/sites-enabled/
   sudo systemctl reload nginx
   ```

### Полезные команды:

```bash
# Просмотр логов
docker-compose logs -f

# Остановка приложения
docker-compose down

# Перезапуск
docker-compose restart

# Обновление
git pull && ./deploy.sh

# Резервное копирование
./backup.sh
```

---

## 📁 Структура проекта

```
ticketGo/
├── cmd/app/           # Точка входа приложения
├── config/            # Конфигурация
├── controllers/       # HTTP контроллеры
├── db/               # Подключение к БД
├── middleware/       # Middleware
├── models/           # Модели данных
├── routes/           # Маршруты
├── service/          # Бизнес-логика
├── static/           # Статические файлы
├── storage/          # Работа с БД
├── templates/        # HTML шаблоны
├── utils/            # Утилиты
├── Dockerfile        # Docker образ
├── docker-compose.yml # Docker Compose
├── deploy.sh         # Скрипт деплоя
├── backup.sh         # Скрипт бэкапа
└── DEPLOYMENT.md     # Руководство по деплою
```

---

## 🔧 Разработка

### Требования:
- Go 1.24+
- PostgreSQL 12+
- Docker (опционально)

### Установка зависимостей:
```bash
go mod download
```

### Запуск тестов:
```bash
go test ./...
```

### Сборка:
```bash
go build -o ticketgo cmd/app/main.go
```

---

## 📞 Поддержка

При возникновении проблем:
1. Проверьте логи: `docker-compose logs`
2. Убедитесь, что все порты открыты
3. Проверьте настройки в `.env` файле
4. Обратитесь к [DEPLOYMENT.md](DEPLOYMENT.md) для подробной диагностики

---
