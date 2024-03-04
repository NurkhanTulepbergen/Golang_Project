# Golang_Project

# JANA Shop

The JANA Shop project is a web application for managing the store menu.

## Description

#######

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
    curl -X POST -H "Content-Type: application/json" -d '{"username":"newuser", "email":"newuser@example.com", "password":"securepassword", "address":"123 Main St"}' http://localhost:2003/users
    ```
## Add Shop

- **Endpoint:**
  - `POST /addshop`

- **Description:**
  - This endpoint allows you to add a new shop to the database.

- **Response:**
  - Status Code: 201 Created
  - Body: JSON object containing the newly created shop's information.
    
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"title":"Example Shop","description":"This is a description of Example Shop."}' http://localhost:2003/addshop
    ```
    
## Request Body

The request body should be a JSON object with the following fields:

- **title**: The title of the shop (required).
- **description**: A description of the shop (optional).

Example:

```json
{
    "title": "Example Shop",
    "description": "This is a description of Example Shop."
}
```



## Get Shop by ID

- **Endpoint:**
  - `GET /shop/:id`

- **Description:**
  - Retrieve information about a shop by providing its ID.

- **Response:**
  - Status Code: 200 OK
  - Body: JSON object containing menu information.

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"title":"Example Shop","description":"This is a description of Example Shop."}' http://localhost:2003/shop/1

    ```

## Update Existing Shop

- **Endpoint:**
  - `PUT /shop/:id`

- **Description:**
  - Update an existing shop.

- **Request Body:**
  - JSON object containing updated shop information.

    ```json
    {
      "name": "Updated Shop",
      "items": ["Item 1", "Item 2", "Item 3"],
    }
    ```

- **Response:**
  - Status Code: 200 OK
  - Body: Success message.

    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Shop", "items":["Item 1", "Item 2", "Item 3"]}' http://localhost:2003/shop/1
    ```

## Delete Shop

- **Endpoint:**
  - `DELETE /shop/:id`

- **Description:**
  - Delete a shop by providing its ID.

- **Response:**
  - Status Code: 200 OK
  - Body: Success message.

    ```bash
    curl -X DELETE http://localhost:2003/shop/1
    ```

---

This documentation provides details on the API routes for creating a new user, getting shop information by ID, updating an existing shop, and deleting a shop. Adjust the examples and details based on your specific project requirements.

  
## Contribution

We appreciate contribution! If you would like to contribute to the project, please create a new branch and submit a pull request.

