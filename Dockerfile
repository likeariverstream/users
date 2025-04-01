FROM golang:1.23 AS builder
ARG ENV
ENV PATH="/go/bin:${PATH}"
WORKDIR /app
COPY . .

RUN go mod download
RUN echo "$ENV" > .env
RUN swag init -g main.go
RUN CGO_ENABLED=0 go build -o ./goapp ./

FROM alpine:latest
ARG APP_TAG
ENV APP_TAG=${APP_TAG}

COPY --from=builder app/goapp /goapp
COPY --from=builder app/.env /.env
COPY --from=builder app/init.sh app/init.sh
COPY --from=builder app/migrations /migrations

ADD https://github.com/pressly/goose/releases/download/v3.24.0/goose_linux_x86_64 /bin/goose

RUN chmod +x /bin/goose

CMD ["sh", "/app/init.sh"]
