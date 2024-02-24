# Golang_Project

# JANA Shop

The JANA Shop project is a web application for managing the store menu.

## Описание

JANA Shop предоставляет следующие возможности:

- Создание новых меню
- Получение информации о конкретном меню по его идентификатору
- Обновление существующего меню
- Удаление меню из системы

## Usage

To run the project, follow these steps:

1. Install the necessary dependencies:

    ```bash
    go mod tidy
    ```

2. Run the project:

    ```bash
    go run main.go
    ```

This will run the project on a local server.

## Маршруты API

- `POST /menus`: Создание нового меню
- `GET /menus/:id`: Получение информации о меню по его идентификатору
- `PUT /menus/:id`: Обновление существующего меню
- `DELETE /menus/:id`: Удаление меню

## Contribution

We appreciate contribution! If you would like to contribute to the project, please create a new branch and submit a pull request.

