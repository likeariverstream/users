FROM golang:1.23 AS builder
LABEL version = 1.0
ARG ENV
ENV PATH="/go/bin:${PATH}"
WORKDIR /app
COPY . .

RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN swag init
RUN go test -v
RUN echo "$ENV" > .env
RUN CGO_ENABLED=0 go build -o userssrv main.go

FROM alpine:latest

COPY --from=builder app/userssrv /userssrv
COPY --from=builder app/.env /.env
COPY --from=builder app/init.sh app/init.sh
COPY --from=builder app/migrations /migrations

ADD https://github.com/pressly/goose/releases/download/v3.24.0/goose_linux_x86_64 /bin/goose

RUN chmod +x /bin/goose

CMD ["sh", "/app/init.sh"]
