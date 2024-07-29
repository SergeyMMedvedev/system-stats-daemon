# Базовый образ Ubuntu
FROM golang:1.22

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Сборка Go-приложения
RUN go mod tidy
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    sysstat sudo net-tools iproute2 tcpdump \
    && rm -rf /var/lib/apt/lists/*

# Определяем команду для запуска приложения
CMD ["bash"]