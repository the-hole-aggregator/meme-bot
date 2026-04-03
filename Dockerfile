FROM golang:1.26.1-alpine

RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow

WORKDIR /app

COPY bot .


CMD [ "./bot" ]
