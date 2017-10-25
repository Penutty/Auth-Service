# Use an official Python runtime as a parent image
FROM golang

RUN apt-get update && apt-get install -y \
    libsybdb5 \
    freetds-dev \
    freetds-common

# Install any needed package dependencies 
RUN go get -u github.com/labstack/echo
RUN go get -u github.com/dgrijalva/jwt-go
RUN go get -u github.com/minus5/gofreetds


# Copy go packages into container.
COPY . /go/src/github.com/penutty/authservice
RUN git clone https://github.com/penutty/dba /go/src/github.com/penutty/dba
RUN git clone https://github.com/penutty/util /go/src/github.com/penutty/util

RUN ls /go/src/github.com/penutty
# Install go packages
RUN go install github.com/penutty/authservice

# Run the outyet command by default when the container starts.
ENTRYPOINT $GOPATH/bin/authservice

# Make port 80 available to the world outside this container
EXPOSE 8080

# Define environment variable
ENV DatabaseConnStr="server=192.168.1.2:1433;database=moment-db;user id=reader;password=123"

# Run app.py when the container launches
CMD "authservice"
