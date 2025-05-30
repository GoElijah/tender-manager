# Сервис для работы с тендерами (MVP)

## Для локального запуска необходимо:

1. Склонировать проект из репозитория на компьютер:
    ```sh
    git clone git@github.com:YanTheBoy/tender-manager.git
    ```

2. Переименовать файл `.env.example` в `.env`:
    ```sh
    mv .env.example .env
    ```

3. В терминале выполнить команду для поднятия докер-контейнера и запуска сервиса на [http://localhost:8080/](http://localhost:8080/):
    ```sh
    make up
    ```

4. В терминале выполнить команду для поднятия локальной базы данных:
    ```sh
    make run-migrate-development
    ```

5. Проверить, что запрос вернул код 200 и сообщение `{"ping":"ok"}`:
    ```sh
    curl localhost:8080/api/ping -i
    ```

---

**Примечание:** Убедитесь, что установлен Docker и Docker Compose, а также утилита sql-migrate для выполнения миграций.