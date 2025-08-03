ARG APP_DIR="/rss-tg-bot"
ARG DIST_DIR=${APP_DIR}/dist/
ARG OUTPUT_FILE_NAME=app

FROM golang:1.25rc2 AS build

ARG APP_DIR
ARG DIST_DIR
ARG OUTPUT_FILE_NAME

WORKDIR ${APP_DIR}

COPY go.mod go.sum ./

RUN go mod download

COPY src/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $OUTPUT_FILE_NAME

FROM alpine:3.21 AS runtime

ARG APP_DIR
ARG DIST_DIR
ARG OUTPUT_FILE_NAME

WORKDIR ${APP_DIR}

RUN apk add --no-cache busybox-suid

COPY crontab.txt /etc/crontabs/root
RUN chmod 0644 /etc/crontabs/root
RUN touch /var/log/cron.log

COPY --from=build ${APP_DIR}/${OUTPUT_FILE_NAME} ${DIST_DIR}/${OUTPUT_FILE_NAME}

CMD ["${DIST_DIR}/app"]
# CMD ["crond", "-f", "-d", "8"]