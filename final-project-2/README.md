# Book List

The simple CRUD application called Book List keeps record of books read. Book list uses the `net/http` package to create simple back end APIs and http client. This project uses JSON for file-based storage, and the CLI for the client. It also has logging capabilities to keep track of the processes as the user interacts with the app. It uses JWT tokens to secure API requests and to keep track of sessions. Finally, the back end is deployed using Docker.

## Project Structure

```
> final-project-2
    > client
        > models                    - data models to be displayed on the CLI
        app.go
    > server
        > config                    - functions for authentication, database connection, storage options
        > controller                - main logic for requests
        > data                      - output location for file-based storage option
        > logs                      - output location for logs
        > models
        > routes                    - routing for the controllers
        > utils                     - utility functions such as better error handling messages
        .env
        app.go
        Dockerfile
    README.md
```
