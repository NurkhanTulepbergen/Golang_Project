# Golang_Project

# JANA Shop
JANA Shop is final project for discipline Developing in Golang

## Team members

| Student name          | Student ID      |
|-----------------------|-----------------|
| Suleimenova Zhasmin   | 22B030444       |
| Taimas Ayazhan        | 22B030447       |
| Taubaev Azamat        | 22B030450       |
| Tulepbergen Nurkhan   | 22B030455       |

## Description
JANA Shop is an online marketplace with categorized product listings where customers can discover a wide range of items, while sellers have the capability to add, remove, and modify the prices of their products. This platform serves as a central hub where buyers can explore diverse product categories, from electronics to fashion, and everything in between. With intuitive navigation and search functionalities, customers can easily find the items they are looking for or explore new offerings. Meanwhile, sellers benefit from the flexibility to manage their inventory, adjust pricing based on market demands, and showcase their products to a broad audience of potential buyers. JANA Shop fosters a dynamic ecosystem where transactions between buyers and sellers are facilitated seamlessly, fostering a vibrant online marketplace experience for all parties involved.

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
# Marketplace REST API 
```
# list of shop function
GET /shop
POST /shop
GET /shop/{id}
PUT /shop/{id}
DELETE /shop/{id}
GET /shop/{shop_id}/product

# list of product fiction
GET /product
POST /product
DELETE /product/{id}
PUT /product/{id}
GET /product/{id}

# list of user function
POST /user
PUT /user/activated
POST /tokens/authentication

# list of order function
GET /order
GET /order/{order_id}
POST /order/{order_id}

# list of history function
GET /history/{userID}
POST /history
DELETE /history/{userID}
PUT /history

# list of follow function
GET /follow/user/{user_id}
POST /follow
DELETE /follow/user/{user_id}/product/{product_id}
```

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
  - `POST /shop`

- **Description:**
  - This endpoint allows you to add a new shop to the database.

- **Response:**
  - Status Code: 201 Created
  - Body: JSON object containing the newly created shop's information.
    
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"title":"Example Shop","description":"This is a description of Example Shop."}' http://localhost:2003/shop
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

## Get Shops

- **Endpoint:**
  - `GET /shops`

- **Description:**
  - Retrieve a list of shops from the database with optional filtering, sorting, and pagination.

- **Query Parameters:**
  - `type`: (Optional) Filter shops by type. Example: `type=electronics`.
  - `sortBy`: (Optional) Sort shops by a specific field. Example: `sortBy=title`.
  - `page`: (Optional) Specify the page number for pagination. Default is 1.
  - `pageSize`: (Optional) Specify the number of shops per page. Default is 10.

- **Response:**
  - Status Code: 200 OK
  - Body: JSON array containing a list of shops based on the provided parameters.

    Example Response Body:

    ```json
    [
        {"id": 1, "title": "Electronics Shop 1", "type": "electronics"},
        {"id": 2, "title": "Electronics Shop 2", "type": "electronics"},
        ...
    ]
    ```

- **Request Example:**
  ```bash
  curl -X GET 'http://your-server-url/shops?type=electronics&sortBy=title&page=1&pageSize=10'
## Filtering:

Use the `type` parameter to filter shops by type. For example, `type=electronics` will only return shops of type "electronics".

## Sorting:

Use the `sortBy` parameter to sort the list of shops by a specific field. For example, `sortBy=title` will sort shops alphabetically by title.

## Pagination:

Use the `page` parameter to specify the page number and the `pageSize` parameter to specify the number of shops per page.

This documentation provides details on the API endpoint for fetching shops with support for filtering by type, sorting, and pagination. Adjust the examples and details based on your specific project requirements.

## Contribution:

We appreciate contributions! If you would like to contribute to the project, please create a new branch and submit a pull request.

This README file explains how to use the `/shops` endpoint with the provided query parameters for filtering by type, sorting, and pagination. Adjustments can be made based on your specific backend implementation and requirements.

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
## Add Product

- **Endpoint:**
  - `POST /product`

- **Description:**
  - This endpoint allows you to add a new product to the database.

- **Request Body:**
  - The request body should be a JSON object with the following fields:
    - **name**: The name of the product (required).
    - **price**: The price of the product (required).
    - **description**: A description of the product (optional).

    Example:

    ```json
    {
        "name": "Example Product",
        "price": 29.99,
        "description": "This is a description of Example Product."
    }
    ```

- **Response:**
  - Status Code: 201 Created
  - Body: JSON object containing the newly created product's information.

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"name":"Example Product","price":29.99,"description":"This is a description of Example Product."}' http://localhost:2003/product
    ```

## Get Product by ID

- **Endpoint:**
  - `GET /product/:id`

- **Description:**
  - Retrieve information about a product by providing its ID.

- **Response:**
  - Status Code: 200 OK
  - Body: JSON object containing product information.

    ```bash
    curl -X GET http://localhost:2003/product/1
    ```

## Update Existing Product

- **Endpoint:**
  - `PUT /product/:id`

- **Description:**
  - Update an existing product.

- **Request Body:**
  - JSON object containing updated product information.

    ```json
    {
        "name": "Updated Product",
        "price": 39.99,
        "description": "Updated description of the product."
    }
    ```

- **Response:**
  - Status Code: 200 OK
  - Body: Success message.

    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Product","price":39.99,"description":"Updated description of the product."}' http://localhost:2003/product/1
    ```

## Delete Product

- **Endpoint:**
  - `DELETE /product/:id`

- **Description:**
  - Delete a product by providing its ID.

- **Response:**
  - Status Code: 200 OK
  - Body: Success message.

    ```bash
    curl -X DELETE http://localhost:2003/product/1
    ```

---
This documentation provides details on the API routes for creating a new user, getting shop information by ID, updating an existing shop, and deleting a shop. Adjust the examples and details based on your specific project requirements.

  
## Contribution

We appreciate contribution! If you would like to contribute to the project, please create a new branch and submit a pull request.

