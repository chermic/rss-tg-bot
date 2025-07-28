Telegram-бот, который парсит RSS-ленты и присылает новости в чат. Скрипт запускается по крон-таске каждый день в 20:00мск

## Настройка GitHub Actions

Для автоматического деплоя через GitHub Actions необходимо настроить следующие секреты в репозитории:

### Необходимые секреты

1. **DOCKERHUB_USERNAME** - имя пользователя в Docker Hub
2. **DOCKERHUB_PASSWORD** - пароль или токен доступа к Docker Hub
3. **TELEGRAM_BOT_TOKEN** - токен Telegram бота
4. **PRODUCTION_SERVER_ADDRESS** - IP адрес или домен продакшн сервера
5. **PRODUCTION_SERVER_USERNAME** - имя пользователя для SSH подключения
6. **PRODUCTION_SERVER_PASSWORD** - пароль для SSH подключения
7. **PRODUCTION_SERVER_PORT** - порт для SSH подключения (обычно 22)

### Настройка секретов

1. Перейдите в настройки репозитория на GitHub
2. Выберите "Secrets and variables" → "Actions"
3. Нажмите "New repository secret"
4. Добавьте каждый из секретов выше

### Автоматический деплой

После настройки секретов, при каждом push в ветку `master` будет автоматически:

1. Собран Docker образ
2. Загружен в Docker Hub
3. Развернут на продакшн сервере