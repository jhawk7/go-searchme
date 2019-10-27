# inherit from the Go official Image
FROM golang:latest

# set a workdir inside docker
WORKDIR /go/src/github.com/dev/go-searchme

# copy . (all in the current directory) to . (WORKDIR)
COPY . .

# run a command - this will run when building the image
RUN go get
RUN go build -o go-searchme

# the port we wish to expose
EXPOSE 8888

# run a command when running the container
CMD ./go-searchme
