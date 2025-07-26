ARG APP_DIR="/app"
ARG OUTPUT_FILE_NAME=main

FROM golang:1.25rc2 AS build

ARG APP_DIR
ARG OUTPUT_FILE_NAME

WORKDIR ${APP_DIR}

COPY go.mod go.sum ./

RUN go mod download

COPY src/*.go .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $OUTPUT_FILE_NAME

FROM alpine:3.21 AS runtime

ARG APP_DIR
ARG OUTPUT_FILE_NAME

WORKDIR ${APP_DIR}

RUN apk add --no-cache busybox-suid

COPY docker/crontab.txt /etc/crontabs/root
RUN chmod 0644 /etc/crontabs/root
RUN touch /var/log/cron.log

COPY --from=build ${APP_DIR}/${OUTPUT_FILE_NAME} .

CMD ["crond", "-f", "-d", "8"]

# CMD ["./app"]