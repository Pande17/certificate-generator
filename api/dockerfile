# use Go for base image
FROM golang:1.22

# set working directory in container
WORKDIR /cmd

# copy go.mod and go.sum, then install dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy all code file to container
COPY . .

# build Go Application
RUN go build -o main ./cmd

# command to run program when this container started
CMD ["./cmd/main.go"]

# export port
EXPOSE 3000