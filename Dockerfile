FROM alpine:3.7

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -g app app \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

COPY . .

RUN chown -R app:app ./

EXPOSE 8888
EXPOSE 9999

USER app

ENTRYPOINT ["/app/mars", "server"]
