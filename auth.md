# Authentication flow

## Registration flow outline

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
 
![Mermaid Diagram](![mermaid-flow](https://github.com/NurkhanTulepbergen/Golang_Project/assets/123255704/b98462e5-c8bf-4a5e-a801-30e45c09bfe1)
)
