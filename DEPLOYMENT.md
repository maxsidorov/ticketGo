# 🚀 Руководство по развертыванию ticketGo на VPS/VDS

## 📋 Требования к серверу

- **ОС**: Ubuntu 20.04+ / CentOS 8+ / Debian 11+
- **RAM**: минимум 2GB (рекомендуется 4GB+)
- **CPU**: 2 ядра (рекомендуется 4+)
- **Диск**: минимум 20GB свободного места
- **Сеть**: статический IP-адрес

## 🔧 Подготовка сервера

### 1. Обновление системы

```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# CentOS/RHEL
sudo yum update -y
```

### 2. Установка Docker

```bash
# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Добавление пользователя в группу docker
sudo usermod -aG docker $USER

# Включение автозапуска Docker
sudo systemctl enable docker
sudo systemctl start docker
```

### 3. Установка Docker Compose

```bash
# Установка Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 4. Установка Nginx (опционально)

```bash
# Ubuntu/Debian
sudo apt install nginx -y

# CentOS/RHEL
sudo yum install nginx -y

# Включение автозапуска
sudo systemctl enable nginx
sudo systemctl start nginx
```

## 📦 Развертывание приложения

### 1. Клонирование репозитория

```bash
# Создание директории для проекта
sudo mkdir -p /opt/ticketgo
sudo chown $USER:$USER /opt/ticketgo
cd /opt/ticketgo

# Клонирование репозитория
git clone https://github.com/your-username/ticketGo.git .
```

### 2. Настройка переменных окружения

Создайте файл `.env` с настройками:

```bash
cat > .env << EOF
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=ticketgo
DB_SSLMODE=disable
PORT=8080
EOF
```

### 3. Запуск приложения

```bash
# Сделать скрипт исполняемым
chmod +x deploy.sh

# Запуск деплоя
./deploy.sh
```

### 4. Настройка Nginx (если используется)

```bash
# Копирование конфигурации
sudo cp nginx.conf /etc/nginx/sites-available/ticketgo

# Создание символической ссылки
sudo ln -s /etc/nginx/sites-available/ticketgo /etc/nginx/sites-enabled/

# Удаление дефолтной конфигурации
sudo rm -f /etc/nginx/sites-enabled/default

# Проверка конфигурации
sudo nginx -t

# Перезапуск Nginx
sudo systemctl reload nginx
```

## 🔒 Настройка безопасности

### 1. Настройка файрвола

```bash
# Ubuntu/Debian (ufw)
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 2. Настройка SSL (Let's Encrypt)

```bash
# Установка Certbot
sudo apt install certbot python3-certbot-nginx -y

# Получение SSL сертификата
sudo certbot --nginx -d your-domain.com -d www.your-domain.com
```

## 📊 Мониторинг и логи

### 1. Просмотр логов

```bash
# Логи приложения
docker-compose logs -f app

# Логи базы данных
docker-compose logs -f postgres

# Логи Nginx
sudo tail -f /var/log/nginx/ticketgo_access.log
sudo tail -f /var/log/nginx/ticketgo_error.log
```

### 2. Мониторинг ресурсов

```bash
# Статус контейнеров
docker-compose ps

# Использование ресурсов
docker stats

# Дисковое пространство
df -h
```

## 🔄 Обновление приложения

### 1. Автоматическое обновление

```bash
# Остановка приложения
docker-compose down

# Получение обновлений
git pull origin main

# Пересборка и запуск
./deploy.sh
```

### 2. Откат к предыдущей версии

```bash
# Переключение на предыдущий коммит
git checkout HEAD~1

# Перезапуск
./deploy.sh
```

## 🛠️ Управление сервисом

### 1. Настройка автозапуска

```bash
# Копирование systemd сервиса
sudo cp systemd/ticketgo.service /etc/systemd/system/

# Перезагрузка systemd
sudo systemctl daemon-reload

# Включение автозапуска
sudo systemctl enable ticketgo

# Запуск сервиса
sudo systemctl start ticketgo
```

### 2. Управление сервисом

```bash
# Статус сервиса
sudo systemctl status ticketgo

# Остановка сервиса
sudo systemctl stop ticketgo

# Перезапуск сервиса
sudo systemctl restart ticketgo
```

## 📈 Резервное копирование

### 1. Резервное копирование базы данных

```bash
# Создание бэкапа
docker-compose exec postgres pg_dump -U postgres ticketgo > backup_$(date +%Y%m%d_%H%M%S).sql

# Восстановление из бэкапа
docker-compose exec -T postgres psql -U postgres ticketgo < backup_file.sql
```

### 2. Автоматическое резервное копирование

Создайте скрипт `backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Бэкап базы данных
docker-compose exec postgres pg_dump -U postgres ticketgo > $BACKUP_DIR/db_backup_$DATE.sql

# Бэкап конфигурации
tar -czf $BACKUP_DIR/config_backup_$DATE.tar.gz .env docker-compose.yml

# Удаление старых бэкапов (старше 30 дней)
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
```

Добавьте в crontab:

```bash
# Редактирование crontab
crontab -e

# Добавление строки для ежедневного бэкапа в 2:00
0 2 * * * /opt/ticketgo/backup.sh
```

## 🚨 Устранение неполадок

### 1. Проверка статуса сервисов

```bash
# Статус Docker
sudo systemctl status docker

# Статус контейнеров
docker-compose ps

# Проверка портов
sudo netstat -tlnp | grep :8080
```

### 2. Перезапуск при проблемах

```bash
# Полный перезапуск
docker-compose down
docker system prune -f
./deploy.sh
```

### 3. Проверка логов

```bash
# Логи приложения
docker-compose logs app

# Логи базы данных
docker-compose logs postgres

# Системные логи
sudo journalctl -u docker
```

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи приложения и системы
2. Убедитесь, что все порты открыты
3. Проверьте доступность базы данных
4. Убедитесь, что Docker и Docker Compose работают корректно

## 🔗 Полезные ссылки

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/) 