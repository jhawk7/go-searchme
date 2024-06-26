FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /go/app
COPY . ./
RUN go mod download
RUN mkdir bin
RUN cd cmd/main/ && go build -o ../../bin/go-searchme

FROM node:current-alpine3.18 AS asset-builder
WORKDIR /frontend
COPY --from=builder /go/app/frontend .
RUN npm install
RUN npm run build

FROM golang:1.22.1-alpine3.19
WORKDIR /go
EXPOSE 8888
COPY --from=builder /go/app/bin/go-searchme ./
COPY --from=asset-builder /frontend/dist ./frontend/dist
ENTRYPOINT ["./go-searchme"]
