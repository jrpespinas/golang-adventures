# Book List

The simple CRUD application called Book List keeps record of books read. Book list uses the `net/http` package to create simple back end APIs and http client. This project uses JSON for file-based storage, and the CLI for the client. It also has logging capabilities to keep track of the processes as the user interacts with the app. It uses JWT tokens to secure API requests and to keep track of sessions. Finally, the back end is deployed using Docker.

## Project Structure

```
> final-project-2
    > client
        > models
        app.go
    > server
        > config
        > controller
        > logs
        > models
        > routes
        > utils
        .env
        app.go
        Dockerfile
    README.md
```
