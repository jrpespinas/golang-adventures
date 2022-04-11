# Book List

The simple CRUD application called Book List keeps record of books read. Book list uses the `net/http` package to create simple back end APIs and http client. This project uses JSON for file-based storage, and the CLI for the client. It also has logging capabilities to keep track of the processes as the user interacts with the app. It uses JWT tokens to secure API requests and to keep track of sessions. Finally, the back end is deployed using Docker.

## Project Structure

```
> final-project-2
    > client
        > models                    - models for displaying data on the CLI
        app.go
    > server
        > config                    - functions for authentication, database connection, storage options
        > controller                - main logic for requests
        > database                  - database methods which handles storage for both file-based storage and MongoDB
        > data                      - output location for file-based storage option
        > logs                      - output location for logs
        > models                    - models for handling incoming and processing data
        > routes                    - routing for the controllers
        > utils                     - utility functions such as better error handling messages, and other middlewares
        .env
        app.go
        go.mod
        Dockerfile
    README.md
    docker-compose.yml
```
