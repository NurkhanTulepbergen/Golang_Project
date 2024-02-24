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
# REST API 

## Create a New User

- **Endpoint:**
  - `POST /users`

- **Description:**
  - Create a new user.

- **Request Body:**
  - JSON object containing user information.

    ```json
    {
      "username": "newuser",
      "email": "newuser@example.com",
      "password": "securepassword",
      "address": "123 Main St"
    }
    ```

- **Response:**
  - Status Code: 200 OK
  - Body: Success message.

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"username":"newuser", "email":"newuser@example.com", "password":"securepassword", "address":"123 Main St"}' http://localhost:8080/users
    ```

## Get Menu by ID

- **Endpoint:**
  - `GET /menus/:id`

- **Description:**
  - Retrieve information about a menu by providing its ID.

- **Response:**
  - Status Code: 200 OK
  - Body: JSON object containing menu information.

    ```bash
    curl http://localhost:8080/menus/1
    ```

## Update Existing Menu

- **Endpoint:**
  - `PUT /menus/:id`

- **Description:**
  - Update an existing menu.

- **Request Body:**
  - JSON object containing updated menu information.

    ```json
    {
      "name": "Updated Menu",
      "items": ["Item 1", "Item 2", "Item 3"],
      "price": 29.99
    }
    ```

- **Response:**
  - Status Code: 200 OK
  - Body: Success message.

    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Menu", "items":["Item 1", "Item 2", "Item 3"], "price":29.99}' http://localhost:8080/menus/1
    ```

## Delete Menu

- **Endpoint:**
  - `DELETE /menus/:id`

- **Description:**
  - Delete a menu by providing its ID.

- **Response:**
  - Status Code: 200 OK
  - Body: Success message.

    ```bash
    curl -X DELETE http://localhost:8080/menus/1
    ```

---

This documentation provides details on the API routes for creating a new user, getting menu information by ID, updating an existing menu, and deleting a menu. Adjust the examples and details based on your specific project requirements.

  
## Contribution

We appreciate contribution! If you would like to contribute to the project, please create a new branch and submit a pull request.

