FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY *.go ./
COPY .env ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /pdf-rest-api

ARG APP_PORT
ENV APP_PORT=${APP_PORT}
EXPOSE ${APP_PORT}

CMD ["/pdf-rest-api"]