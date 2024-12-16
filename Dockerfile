FROM golang:1.23.2

# Set the Current Working Directory inside the container
WORKDIR /EcoHero/src

COPY go.mod go.sum ./

# if you prefer to automate the module initialization and dependency installation, you can uncomment the following lines
# Initialize a Go module and get dependencies
# RUN go mod init EcoHero
RUN go mod tidy


# Copy go mod and sum files
# installs gorilla library
# RUN go get -u github.com/gorilla/mux
# RUN go get -u github.com/go-sql-driver/mysql

# Copy the source from the current directory to the Working Directory /go/src inside the container
COPY . .

# Build the Go app
CMD ["go", "run", "main.go"]