# syntax=docker/dockerfile:1
FROM golang:1.17-alpine

# Set the work directory of the image
WORKDIR /app

# Download necessary Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the source code
# include files to ignore in .dockerignore
COPY . .

# Compile the project
RUN go build -o book-list

# Expose the port to the container
# This port is not the same port to access 
# the application through the localhost
EXPOSE 8000

# Run the binary executable file
CMD ["go", "run", "book-list"]