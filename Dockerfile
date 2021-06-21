# inherit from the Go official Image
FROM golang:1.15-alpine3.13

# set a workdir inside docker
WORKDIR /go/src/github.com/dev/go-searchme

# copy . (all in the current directory) to . (WORKDIR)
COPY . .

# run a command - this will run when building the image
RUN mkdir bin
RUN go build -o bin/go-searchme

# the port we wish to expose
EXPOSE 8888

# run a command when running the container
CMD ./bin/go-searchme
