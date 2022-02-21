Book List
===
## Overview
This project represents the culmination of my Golang training. The simple CRUD application called Book List keeps record of books read. Book list uses the `net/http` package to create simple back end APIs and http client. This project uses JSON for file-based storage, and the CLI for the client. It also has logging capabilities to keep track of the processes as the user interacts with the app. It uses JWT tokens to secure API requests and to keep track of sessions. Finally, the back end is deployed using Docker.

## Usage
To run this project on your machine make sure you have Docker and Golang installed. Here is the link for [Golang installation](https://go.dev/doc/install), and here is for [Docker](https://www.docker.com/products/docker-desktop) which you can find right away on the landing page.

If you already have Docker and Golang installed, follow these steps to run the project:

### Server
1. Clone the repository, and change directory to the `final-project`.
```bash
$ git clone git@github.com:jrpespinas/golang-adventures.git
$ cd final-project
```
2. Change directory to `server` then deploy the back end.
```bash
$ cd server
$ docker build -t book-list .
```
3. Once the build finishes, proceed to run the Docker image.
```bash
$ docker run -p 8080:8080 book-list
```
You will notice the logs being displayed on the terminal which means running the back end server is successful.

### Client
1. Open a separate terminal and go to the `final-project` once more, then finally go to the `client` folder.
```bash
$ cd final-project
$ cd client
```
2. Build client
```bash
$ go build -o client
```

3. Run the client
```bash
$ ./client
```
Now from the client, you may now interact with the application. I suggest you choose the **signup** command then observe the logs you have opened from the other terminal. This is for you to create an account and to check if everything is working.

## License
```
MIT License

Copyright (c) 2022 Jan Rodolf Espinas

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
