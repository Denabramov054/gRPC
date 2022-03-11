FROM golang:alpine

RUN go version

ENV GOPATH=/bin

WORKDIR /app
COPY ./ /app

RUN apk update \
    apk add postgresql-client

RUN go mod download
RUN go build -o /bin/gRPCpet ./cmd/server/main.go

EXPOSE ${APP_PORT}

CMD ["/bin/gRPCpet"]