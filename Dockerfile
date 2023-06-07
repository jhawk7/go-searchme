FROM golang:1.20.0-alpine3.17 AS builder
WORKDIR /go/app
COPY . ./
RUN go mod download
RUN mkdir bin
RUN cd cmd/main/ && go build -o ../../bin/go-searchme

FROM golang:1.20.0-alpine3.17
WORKDIR /go
EXPOSE 8888
COPY --from=builder /go/app/bin/go-searchme ./
COPY templates/ templates/
ENTRYPOINT ["./go-searchme"]
