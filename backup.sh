#!/bin/bash

# Скрипт для автоматического резервного копирования ticketGo
# Рекомендуется добавить в crontab для ежедневного выполнения

set -e

# Настройки
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
PROJECT_DIR="/opt/ticketgo"
RETENTION_DAYS=30

# Создание директории для бэкапов
mkdir -p $BACKUP_DIR

echo "🔄 Начинаем резервное копирование ticketGo..."

# Переход в директорию проекта
cd $PROJECT_DIR

# Бэкап базы данных
echo "📊 Создание бэкапа базы данных..."
docker-compose exec -T postgres pg_dump -U postgres ticketgo > $BACKUP_DIR/db_backup_$DATE.sql

# Бэкап конфигурации
echo "⚙️ Создание бэкапа конфигурации..."
tar -czf $BACKUP_DIR/config_backup_$DATE.tar.gz .env docker-compose.yml

# Бэкап кода (опционально)
echo "💻 Создание бэкапа кода..."
tar -czf $BACKUP_DIR/code_backup_$DATE.tar.gz --exclude='.git' --exclude='node_modules' .

# Проверка размера бэкапов
echo "📏 Размер бэкапов:"
ls -lh $BACKUP_DIR/*_$DATE.*

# Удаление старых бэкапов
echo "🧹 Удаление старых бэкапов (старше $RETENTION_DAYS дней)..."
find $BACKUP_DIR -name "*.sql" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete

# Проверка свободного места
echo "💾 Свободное место на диске:"
df -h $BACKUP_DIR

echo "✅ Резервное копирование завершено успешно!"
echo "📁 Бэкапы сохранены в: $BACKUP_DIR" 