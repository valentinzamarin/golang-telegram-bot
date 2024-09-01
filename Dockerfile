# Используем официальный образ Go для сборки
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем все остальные файлы
COPY . .

# Собираем приложение
RUN go build -o bot .

# Используем тот же образ Go для запуска приложения
FROM golang:1.20

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/bot /bot

# Устанавливаем переменные окружения (заполняйте ваши реальные ключи в .env файле)
ENV TELEGRAM_TOKEN=your_telegram_bot_token
ENV WEATHER_API_KEY=your_openweathermap_api_key

# Запускаем приложение
CMD ["/bot"]