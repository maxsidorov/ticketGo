#!/bin/bash

# Скрипт для деплоя ticketGo на VPS/VDS сервер
# Использование: ./deploy.sh [production|staging]

set -e

ENVIRONMENT=${1:-production}
PROJECT_NAME="ticketgo"
DOCKER_IMAGE="ticketgo:latest"

echo "🚀 Начинаем деплой ticketGo в окружении: $ENVIRONMENT"

# Проверяем наличие Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен. Установите Docker и попробуйте снова."
    exit 1
fi

# Проверяем наличие Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен. Установите Docker Compose и попробуйте снова."
    exit 1
fi

echo "📦 Собираем Docker образ..."
docker build -t $DOCKER_IMAGE .

echo "🔄 Останавливаем существующие контейнеры..."
docker-compose down || true

echo "🧹 Очищаем неиспользуемые образы..."
docker image prune -f

echo "🚀 Запускаем приложение..."
docker-compose up -d

echo "⏳ Ждем запуска сервисов..."
sleep 10

# Проверяем статус контейнеров
echo "📊 Статус контейнеров:"
docker-compose ps

# Проверяем логи приложения
echo "📋 Логи приложения:"
docker-compose logs app --tail=20

echo "✅ Деплой завершен успешно!"
echo "🌐 Приложение доступно по адресу: http://localhost:8080"
echo "🗄️  База данных PostgreSQL доступна на порту 5432"

# Показываем команды для управления
echo ""
echo "📝 Полезные команды:"
echo "  Просмотр логов: docker-compose logs -f"
echo "  Остановка: docker-compose down"
echo "  Перезапуск: docker-compose restart"
echo "  Обновление: ./deploy.sh" 