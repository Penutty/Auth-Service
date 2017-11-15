# Use an official Python runtime as a parent image
FROM golang:alpine

RUN apk add --update \
    git \
    freetds-dev \
    g++

# Install any needed package dependencies 
RUN go get -u github.com/dgrijalva/jwt-go
RUN go get -u github.com/minus5/gofreetds
RUN go get -u github.com/Masterminds/squirrel

# Copy go packages into container.
COPY . /go/src/github.com/penutty/authservice

# Create log folder
RUN mkdir /go/log 

# Install go packages
RUN go install github.com/penutty/authservice

# Run the outyet command by default when the container starts.
ENTRYPOINT $GOPATH/bin/authservice

# Make port 80 available to the world outside this container
EXPOSE 8080

# Define environment variable
ENV DatabaseConnStr="server=192.168.1.2:1433;database=moment-db;user id=reader;password=123"

