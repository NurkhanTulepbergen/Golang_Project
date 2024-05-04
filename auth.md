##
```mermaid
sequenceDiagram
    participant User
    participant Server
    participant Database

    User->>Server: createUser(name, password)
    Server->>Database: Open database connection
    Database-->>Server: Database connection established
    Server->>Database: Execute INSERT query
    Database-->>Server: INSERT query executed
    Server-->>User: Print user ID
    end
```

```mermaid
sequenceDiagram
    participant User
    participant Server
    participant Database

    User->>Server: createProduct(title, description, price)
    Server->>Database: Open database connection
    Database-->>Server: Database connection established
    Server->>Database: Execute INSERT query
    Database-->>Server: INSERT query executed
    Server-->>User: Print success message
    end
```
